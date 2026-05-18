package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/service"
	"github.com/go-playground/validator/v10"
)

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// siakadLoginRequest is the request body for SIAKAD login.
type siakadLoginRequest struct {
	NIM      string `json:"nim" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// githubCallbackRequest is the request body for GitHub OAuth callback.
type githubCallbackRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	Code      string `json:"code" validate:"required"`
}

// HandleSIAKADLogin handles POST /api/auth/siakad-login.
func (h *AuthHandler) HandleSIAKADLogin(w http.ResponseWriter, r *http.Request) {
	var req siakadLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondError(w, http.StatusBadRequest, "nim and password are required")
		return
	}

	result, err := h.authService.InitiateSIAKADLogin(r.Context(), req.NIM, req.Password)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"session_id":   result.SessionID,
		"redirect_url": result.RedirectURL,
	})
}

// HandleGitHubCallback handles POST /api/auth/github-callback.
func (h *AuthHandler) HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	var req githubCallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondError(w, http.StatusBadRequest, "session_id and code are required")
		return
	}

	result, err := h.authService.CompleteGitHubOAuth(r.Context(), req.SessionID, req.Code)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	// Set refresh token as httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(7 * 24 * time.Hour / time.Second),
	})

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"user":         result.User,
		"access_token": result.AccessToken,
	})
}

// HandleRefreshToken handles POST /api/auth/refresh.
func (h *AuthHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		RespondError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	result, err := h.authService.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	// Set new refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(7 * 24 * time.Hour / time.Second),
	})

	RespondJSON(w, http.StatusOK, map[string]string{
		"access_token": result.AccessToken,
	})
}

// HandleLogout handles POST /api/auth/logout.
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.authService.Logout(r.Context(), claims.UserID); err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

// handleAuthError maps domain errors to HTTP responses.
func (h *AuthHandler) handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		RespondError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, domain.ErrDeviceRejected):
		RespondError(w, http.StatusForbidden, "device rejected by SIAKAD")
	case errors.Is(err, domain.ErrSIAKADUnavailable):
		RespondError(w, http.StatusServiceUnavailable, "SIAKAD service unavailable")
	case errors.Is(err, domain.ErrSessionInitFailed):
		RespondError(w, http.StatusBadGateway, "session initialization failed")
	case errors.Is(err, domain.ErrTokenInvalid):
		RespondError(w, http.StatusUnauthorized, "invalid or expired session")
	case errors.Is(err, domain.ErrTokenExpired):
		RespondError(w, http.StatusUnauthorized, "token expired")
	case errors.Is(err, domain.ErrForbidden):
		RespondError(w, http.StatusForbidden, "forbidden")
	default:
		RespondError(w, http.StatusInternalServerError, "internal server error")
	}
}
