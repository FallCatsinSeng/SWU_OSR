package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ForumHandler handles forum HTTP requests.
type ForumHandler struct {
	forumService service.ForumService
}

// NewForumHandler creates a new forum handler.
func NewForumHandler(forumService service.ForumService) *ForumHandler {
	return &ForumHandler{
		forumService: forumService,
	}
}

// HandleListThreads handles GET /api/repos/{id}/threads.
func (h *ForumHandler) HandleListThreads(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	repoID, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid repo id")
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	params := domain.PaginationParams{
		Cursor: cursor,
		Limit:  limit,
	}

	result, err := h.forumService.ListThreads(r.Context(), repoID, params)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to list threads")
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

// createThreadRequest is the request body for creating a thread.
type createThreadRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// HandleCreateThread handles POST /api/repos/{id}/threads (auth required).
func (h *ForumHandler) HandleCreateThread(w http.ResponseWriter, r *http.Request) {
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

	var req createThreadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := domain.CreateThreadInput{
		RepoID:   repoID,
		AuthorID: claims.UserID,
		Title:    req.Title,
		Body:     req.Body,
	}

	thread, err := h.forumService.CreateThread(r.Context(), input)
	if err != nil {
		h.handleForumError(w, err)
		return
	}

	RespondJSON(w, http.StatusCreated, thread)
}

// HandleGetThread handles GET /api/threads/{id}.
func (h *ForumHandler) HandleGetThread(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	threadID, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid thread id")
		return
	}

	thread, comments, err := h.forumService.GetThread(r.Context(), threadID)
	if err != nil {
		RespondError(w, http.StatusNotFound, "thread not found")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"thread":   thread,
		"comments": comments,
	})
}

// createCommentRequest is the request body for creating a comment.
type createCommentRequest struct {
	Body     string     `json:"body"`
	ParentID *uuid.UUID `json:"parent_id"`
}

// HandleCreateComment handles POST /api/threads/{id}/comments (auth required).
func (h *ForumHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	threadID, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid thread id")
		return
	}

	var req createCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := domain.CreateCommentInput{
		ThreadID: threadID,
		AuthorID: claims.UserID,
		ParentID: req.ParentID,
		Body:     req.Body,
	}

	comment, err := h.forumService.CreateComment(r.Context(), input)
	if err != nil {
		h.handleForumError(w, err)
		return
	}

	RespondJSON(w, http.StatusCreated, comment)
}

// HandleListNotifications handles GET /api/notifications (auth required).
func (h *ForumHandler) HandleListNotifications(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	notifications, err := h.forumService.ListNotifications(r.Context(), claims.UserID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to list notifications")
		return
	}

	RespondJSON(w, http.StatusOK, notifications)
}

// HandleMarkNotificationRead handles PUT /api/notifications/{id}/read (auth required).
func (h *ForumHandler) HandleMarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	notifID, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	if err := h.forumService.MarkNotificationRead(r.Context(), claims.UserID, notifID); err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to mark notification as read")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "notification marked as read"})
}

// handleForumError maps service errors to appropriate HTTP responses.
// Validation errors (containing known patterns) return 400 with the message.
// ErrNotFound returns 404 with a generic message.
// All other errors return 500 with a generic message to avoid leaking internals.
func (h *ForumHandler) handleForumError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		RespondError(w, http.StatusNotFound, "resource not found")
	case isValidationError(err):
		RespondError(w, http.StatusBadRequest, err.Error())
	default:
		RespondError(w, http.StatusInternalServerError, "internal server error")
	}
}

// isValidationError checks if the error message matches known validation patterns
// from the forum service layer.
func isValidationError(err error) bool {
	msg := err.Error()
	validationPatterns := []string{
		"must be at least",
		"must not exceed",
		"must not be empty",
	}
	for _, pattern := range validationPatterns {
		if strings.Contains(msg, pattern) {
			return true
		}
	}
	return false
}
