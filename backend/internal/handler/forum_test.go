package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockForumService implements service.ForumService for testing.
type mockForumService struct {
	listThreadsFn          func(ctx context.Context, repoID uuid.UUID, params domain.PaginationParams) (*domain.ThreadList, error)
	createThreadFn         func(ctx context.Context, input domain.CreateThreadInput) (*domain.Thread, error)
	getThreadFn            func(ctx context.Context, threadID uuid.UUID) (*domain.Thread, []domain.Comment, error)
	createCommentFn        func(ctx context.Context, input domain.CreateCommentInput) (*domain.Comment, error)
	listNotificationsFn    func(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error)
	markNotificationReadFn func(ctx context.Context, userID, notifID uuid.UUID) error
}

func (m *mockForumService) ListThreads(ctx context.Context, repoID uuid.UUID, params domain.PaginationParams) (*domain.ThreadList, error) {
	return m.listThreadsFn(ctx, repoID, params)
}
func (m *mockForumService) CreateThread(ctx context.Context, input domain.CreateThreadInput) (*domain.Thread, error) {
	return m.createThreadFn(ctx, input)
}
func (m *mockForumService) GetThread(ctx context.Context, threadID uuid.UUID) (*domain.Thread, []domain.Comment, error) {
	return m.getThreadFn(ctx, threadID)
}
func (m *mockForumService) CreateComment(ctx context.Context, input domain.CreateCommentInput) (*domain.Comment, error) {
	return m.createCommentFn(ctx, input)
}
func (m *mockForumService) ListNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error) {
	return m.listNotificationsFn(ctx, userID)
}
func (m *mockForumService) MarkNotificationRead(ctx context.Context, userID, notifID uuid.UUID) error {
	return m.markNotificationReadFn(ctx, userID, notifID)
}

// --- HandleListThreads ---

func TestHandleListThreads_Success(t *testing.T) {
	repoID := uuid.New()
	svc := &mockForumService{
		listThreadsFn: func(_ context.Context, rid uuid.UUID, _ domain.PaginationParams) (*domain.ThreadList, error) {
			assert.Equal(t, repoID, rid)
			return &domain.ThreadList{Threads: []domain.Thread{}}, nil
		},
	}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/repos/"+repoID.String()+"/threads", nil)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleListThreads(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleListThreads_InvalidRepoID(t *testing.T) {
	svc := &mockForumService{}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/repos/bad/threads", nil)
	req = injectChiParam(req, "id", "bad")
	w := httptest.NewRecorder()

	h.HandleListThreads(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleListThreads_WithCustomLimit(t *testing.T) {
	repoID := uuid.New()
	var gotLimit int
	svc := &mockForumService{
		listThreadsFn: func(_ context.Context, _ uuid.UUID, params domain.PaginationParams) (*domain.ThreadList, error) {
			gotLimit = params.Limit
			return &domain.ThreadList{Threads: []domain.Thread{}}, nil
		},
	}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/repos/"+repoID.String()+"/threads?limit=5", nil)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleListThreads(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 5, gotLimit)
}

// --- HandleCreateThread ---

func TestHandleCreateThread_Success(t *testing.T) {
	userID := uuid.New()
	repoID := uuid.New()
	threadID := uuid.New()
	svc := &mockForumService{
		createThreadFn: func(_ context.Context, input domain.CreateThreadInput) (*domain.Thread, error) {
			assert.Equal(t, userID, input.AuthorID)
			assert.Equal(t, repoID, input.RepoID)
			return &domain.Thread{ID: threadID, Title: input.Title}, nil
		},
	}
	h := NewForumHandler(svc)

	body := `{"title":"Test Thread Title","body":"This is the body text for the thread."}`
	req := httptest.NewRequest(http.MethodPost, "/api/repos/"+repoID.String()+"/threads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleCreateThread(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandleCreateThread_Unauthorized(t *testing.T) {
	svc := &mockForumService{}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/repos/"+uuid.New().String()+"/threads", strings.NewReader(`{}`))
	req = injectChiParam(req, "id", uuid.New().String())
	w := httptest.NewRecorder()

	h.HandleCreateThread(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleCreateThread_InvalidJSON(t *testing.T) {
	userID := uuid.New()
	repoID := uuid.New()
	svc := &mockForumService{}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/repos/"+repoID.String()+"/threads", strings.NewReader("not json"))
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleCreateThread(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreateThread_ValidationError(t *testing.T) {
	userID := uuid.New()
	repoID := uuid.New()
	svc := &mockForumService{
		createThreadFn: func(_ context.Context, _ domain.CreateThreadInput) (*domain.Thread, error) {
			return nil, fmt.Errorf("title must be at least 5 characters")
		},
	}
	h := NewForumHandler(svc)

	body := `{"title":"Hi","body":"body"}`
	req := httptest.NewRequest(http.MethodPost, "/api/repos/"+repoID.String()+"/threads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleCreateThread(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- HandleGetThread ---

func TestHandleGetThread_Success(t *testing.T) {
	threadID := uuid.New()
	svc := &mockForumService{
		getThreadFn: func(_ context.Context, id uuid.UUID) (*domain.Thread, []domain.Comment, error) {
			return &domain.Thread{ID: id, Title: "Test"}, []domain.Comment{}, nil
		},
	}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/threads/"+threadID.String(), nil)
	req = injectChiParam(req, "id", threadID.String())
	w := httptest.NewRecorder()

	h.HandleGetThread(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body envelope
	err := json.NewDecoder(w.Body).Decode(&body)
	require.NoError(t, err)
	assert.True(t, body.OK)

	data, ok := body.Data.(map[string]interface{})
	require.True(t, ok, "data should be a map")
	assert.Contains(t, data, "thread")
	assert.Contains(t, data, "comments")
}

func TestHandleGetThread_InvalidID(t *testing.T) {
	svc := &mockForumService{}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/threads/bad", nil)
	req = injectChiParam(req, "id", "bad")
	w := httptest.NewRecorder()

	h.HandleGetThread(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleGetThread_NotFound(t *testing.T) {
	threadID := uuid.New()
	svc := &mockForumService{
		getThreadFn: func(_ context.Context, _ uuid.UUID) (*domain.Thread, []domain.Comment, error) {
			return nil, nil, domain.ErrNotFound
		},
	}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/threads/"+threadID.String(), nil)
	req = injectChiParam(req, "id", threadID.String())
	w := httptest.NewRecorder()

	h.HandleGetThread(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --- HandleCreateComment ---

func TestHandleCreateComment_Success(t *testing.T) {
	userID := uuid.New()
	threadID := uuid.New()
	commentID := uuid.New()
	svc := &mockForumService{
		createCommentFn: func(_ context.Context, input domain.CreateCommentInput) (*domain.Comment, error) {
			assert.Equal(t, userID, input.AuthorID)
			assert.Equal(t, threadID, input.ThreadID)
			return &domain.Comment{ID: commentID, Body: input.Body}, nil
		},
	}
	h := NewForumHandler(svc)

	body := `{"body":"This is a comment."}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/"+threadID.String()+"/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", threadID.String())
	w := httptest.NewRecorder()

	h.HandleCreateComment(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandleCreateComment_Unauthorized(t *testing.T) {
	svc := &mockForumService{}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/threads/"+uuid.New().String()+"/comments", strings.NewReader(`{}`))
	req = injectChiParam(req, "id", uuid.New().String())
	w := httptest.NewRecorder()

	h.HandleCreateComment(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleCreateComment_ValidationError(t *testing.T) {
	userID := uuid.New()
	threadID := uuid.New()
	svc := &mockForumService{
		createCommentFn: func(_ context.Context, _ domain.CreateCommentInput) (*domain.Comment, error) {
			return nil, fmt.Errorf("body must not be empty")
		},
	}
	h := NewForumHandler(svc)

	body := `{"body":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/"+threadID.String()+"/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", threadID.String())
	w := httptest.NewRecorder()

	h.HandleCreateComment(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- HandleListNotifications ---

func TestHandleListNotifications_Success(t *testing.T) {
	userID := uuid.New()
	svc := &mockForumService{
		listNotificationsFn: func(_ context.Context, uid uuid.UUID) ([]domain.Notification, error) {
			assert.Equal(t, userID, uid)
			return []domain.Notification{
				{ID: uuid.New(), UserID: userID, Type: "new_thread", Message: "hello"},
			}, nil
		},
	}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	req = injectClaims(req, userID)
	w := httptest.NewRecorder()

	h.HandleListNotifications(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleListNotifications_Unauthorized(t *testing.T) {
	svc := &mockForumService{}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	w := httptest.NewRecorder()

	h.HandleListNotifications(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- HandleMarkNotificationRead ---

func TestHandleMarkNotificationRead_Success(t *testing.T) {
	userID := uuid.New()
	notifID := uuid.New()
	svc := &mockForumService{
		markNotificationReadFn: func(_ context.Context, uid, nid uuid.UUID) error {
			assert.Equal(t, userID, uid)
			assert.Equal(t, notifID, nid)
			return nil
		},
	}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/notifications/"+notifID.String()+"/read", nil)
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", notifID.String())
	w := httptest.NewRecorder()

	h.HandleMarkNotificationRead(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleMarkNotificationRead_InvalidID(t *testing.T) {
	userID := uuid.New()
	svc := &mockForumService{}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/notifications/bad/read", nil)
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", "bad")
	w := httptest.NewRecorder()

	h.HandleMarkNotificationRead(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleMarkNotificationRead_Unauthorized(t *testing.T) {
	svc := &mockForumService{}
	h := NewForumHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/notifications/"+uuid.New().String()+"/read", nil)
	req = injectChiParam(req, "id", uuid.New().String())
	w := httptest.NewRecorder()

	h.HandleMarkNotificationRead(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
