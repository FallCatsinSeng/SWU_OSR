package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// ProfileHandler handles profile HTTP requests.
type ProfileHandler struct {
	profileService service.ProfileService
	validate       *validator.Validate
}

// NewProfileHandler creates a new profile handler.
func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		validate:       validator.New(),
	}
}

// HandleGetPublicProfile handles GET /api/profiles/{alias}.
func (h *ProfileHandler) HandleGetPublicProfile(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		RespondError(w, http.StatusBadRequest, "alias is required")
		return
	}

	profile, err := h.profileService.GetPublicProfile(r.Context(), alias)
	if err != nil {
		h.handleProfileError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, profile)
}

// updateProfileRequest is the request body for updating a profile.
type updateProfileRequest struct {
	Alias     string `json:"alias"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

// HandleUpdateProfile handles PUT /api/profile (auth required).
func (h *ProfileHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := service.UpdateProfileInput{
		Alias:     req.Alias,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
	}

	if err := h.profileService.UpdateProfile(r.Context(), claims.UserID, input); err != nil {
		h.handleProfileError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "profile updated"})
}

// HandleGetRealIdentity handles GET /api/profiles/{alias}/identity (auth required).
func (h *ProfileHandler) HandleGetRealIdentity(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	alias := chi.URLParam(r, "alias")
	if alias == "" {
		RespondError(w, http.StatusBadRequest, "alias is required")
		return
	}

	identity, err := h.profileService.GetRealIdentity(r.Context(), claims.UserID, alias)
	if err != nil {
		h.handleProfileError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, identity)
}

// handleProfileError maps domain errors to HTTP responses.
func (h *ProfileHandler) handleProfileError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		RespondError(w, http.StatusNotFound, "not found")
	case errors.Is(err, domain.ErrForbidden):
		RespondError(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, domain.ErrDuplicateAlias):
		RespondError(w, http.StatusConflict, "alias already taken or invalid format")
	default:
		RespondError(w, http.StatusInternalServerError, "internal server error")
	}
}
