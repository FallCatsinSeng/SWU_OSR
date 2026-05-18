package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ShowcaseHandler handles showcase HTTP requests.
type ShowcaseHandler struct {
	showcaseService service.ShowcaseService
	validate        *validator.Validate
}

// NewShowcaseHandler creates a new showcase handler.
func NewShowcaseHandler(showcaseService service.ShowcaseService) *ShowcaseHandler {
	return &ShowcaseHandler{
		showcaseService: showcaseService,
		validate:        validator.New(),
	}
}

// HandleGetAvailableRepos handles GET /api/repos/available (auth required).
func (h *ShowcaseHandler) HandleGetAvailableRepos(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	repos, err := h.showcaseService.GetAvailableRepos(r.Context(), claims.UserID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to fetch repositories")
		return
	}

	RespondJSON(w, http.StatusOK, repos)
}

// setShowcaseRequest is the request body for setting showcase repos.
type setShowcaseRequest struct {
	Selections []domain.ShowcaseSelection `json:"selections" validate:"required,max=20,dive"`
}

// HandleSetShowcase handles POST /api/showcase (auth required).
func (h *ShowcaseHandler) HandleSetShowcase(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req setShowcaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondError(w, http.StatusBadRequest, "validation failed: selections required, max 20")
		return
	}

	if err := h.showcaseService.SetShowcase(r.Context(), claims.UserID, req.Selections); err != nil {
		h.handleShowcaseError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "showcase updated"})
}

// HandleGetShowcase handles GET /api/showcase (auth required).
func (h *ShowcaseHandler) HandleGetShowcase(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	repos, err := h.showcaseService.GetShowcase(r.Context(), claims.UserID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to fetch showcase")
		return
	}

	RespondJSON(w, http.StatusOK, repos)
}

// HandleRemoveFromShowcase handles DELETE /api/showcase/{id} (auth required).
func (h *ShowcaseHandler) HandleRemoveFromShowcase(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	repoID, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid repo id")
		return
	}

	if err := h.showcaseService.RemoveFromShowcase(r.Context(), claims.UserID, repoID); err != nil {
		h.handleShowcaseError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "repo removed from showcase"})
}

// handleShowcaseError maps domain errors to HTTP responses.
func (h *ShowcaseHandler) handleShowcaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		RespondError(w, http.StatusNotFound, "not found")
	default:
		RespondError(w, http.StatusInternalServerError, "internal server error")
	}
}
