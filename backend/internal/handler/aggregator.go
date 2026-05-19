package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// AggregatorHandler handles aggregator HTTP requests.
type AggregatorHandler struct {
	aggregatorService service.AggregatorService
}

// NewAggregatorHandler creates a new aggregator handler.
func NewAggregatorHandler(aggregatorService service.AggregatorService) *AggregatorHandler {
	return &AggregatorHandler{
		aggregatorService: aggregatorService,
	}
}

// HandleWebhook handles POST /api/webhooks/github.
// Always returns 200 to GitHub to prevent retry storms.
func (h *AggregatorHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Limit request body to 5MB to prevent memory exhaustion from oversized payloads
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		// Still return 200 to GitHub
		RespondJSON(w, http.StatusOK, map[string]string{"status": "error"})
		return
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	eventType := r.Header.Get("X-GitHub-Event")
	deliveryID := r.Header.Get("X-GitHub-Delivery")

	_ = h.aggregatorService.ProcessWebhook(r.Context(), payload, signature, eventType, deliveryID)

	// Always return 200
	RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// HandleGetFeed handles GET /api/feed.
func (h *AggregatorHandler) HandleGetFeed(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	params := domain.FeedParams{
		Cursor: cursor,
		Limit:  limit,
	}

	result, err := h.aggregatorService.GetActivityFeed(r.Context(), params)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to fetch feed")
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

// HandleSyncActivity handles POST /api/activity/sync (auth required).
// It fetches recent public GitHub events for the current user and inserts new ones into the feed.
func (h *AggregatorHandler) HandleSyncActivity(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	inserted, err := h.aggregatorService.SyncUserActivity(r.Context(), claims.UserID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to sync activity")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"synced": inserted,
	})
}

// HandleGetUserActivity handles GET /api/users/{id}/activity.
func (h *AggregatorHandler) HandleGetUserActivity(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	params := domain.FeedParams{
		Cursor: cursor,
		Limit:  limit,
	}

	result, err := h.aggregatorService.GetUserActivity(r.Context(), userID, params)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to fetch user activity")
		return
	}

	RespondJSON(w, http.StatusOK, result)
}
