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
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/github"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockShowcaseService implements service.ShowcaseService for testing.
type mockShowcaseService struct {
	getAvailableReposFn     func(ctx context.Context, userID uuid.UUID) ([]github.Repository, error)
	setShowcaseFn           func(ctx context.Context, userID uuid.UUID, selections []domain.ShowcaseSelection) error
	getShowcaseFn           func(ctx context.Context, userID uuid.UUID) ([]domain.ShowcaseRepo, error)
	removeFromShowcaseFn    func(ctx context.Context, userID, repoID uuid.UUID) error
	getPublicRepoFn         func(ctx context.Context, repoID uuid.UUID) (*domain.ShowcaseRepo, error)
	updateRepoDescriptionFn func(ctx context.Context, userID, repoID uuid.UUID, description string) error
}

func (m *mockShowcaseService) GetAvailableRepos(ctx context.Context, userID uuid.UUID) ([]github.Repository, error) {
	return m.getAvailableReposFn(ctx, userID)
}
func (m *mockShowcaseService) SetShowcase(ctx context.Context, userID uuid.UUID, selections []domain.ShowcaseSelection) error {
	return m.setShowcaseFn(ctx, userID, selections)
}
func (m *mockShowcaseService) GetShowcase(ctx context.Context, userID uuid.UUID) ([]domain.ShowcaseRepo, error) {
	return m.getShowcaseFn(ctx, userID)
}
func (m *mockShowcaseService) RemoveFromShowcase(ctx context.Context, userID, repoID uuid.UUID) error {
	return m.removeFromShowcaseFn(ctx, userID, repoID)
}
func (m *mockShowcaseService) GetPublicRepo(ctx context.Context, repoID uuid.UUID) (*domain.ShowcaseRepo, error) {
	return m.getPublicRepoFn(ctx, repoID)
}
func (m *mockShowcaseService) UpdateRepoDescription(ctx context.Context, userID, repoID uuid.UUID, description string) error {
	return m.updateRepoDescriptionFn(ctx, userID, repoID, description)
}
func (m *mockShowcaseService) UpdateShowcase(ctx context.Context, userID uuid.UUID, selections []domain.ShowcaseSelection) error {
	return m.setShowcaseFn(ctx, userID, selections)
}
func (m *mockShowcaseService) SetAggregatorService(_ service.AggregatorService) {} // no-op for tests

// injectClaims injects JWT claims into a request context, simulating the auth middleware.
func injectClaims(r *http.Request, userID uuid.UUID) *http.Request {
	claims := &domain.UserClaims{
		UserID: userID,
		Alias:  "testuser",
		Role:   domain.RoleStudent,
	}
	ctx := context.WithValue(r.Context(), domain.UserClaimsKey, claims)
	return r.WithContext(ctx)
}

// injectChiParam sets a Chi URL parameter in the request context.
func injectChiParam(r *http.Request, key, val string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// --- HandleGetAvailableRepos ---

func TestHandleGetAvailableRepos_Success(t *testing.T) {
	userID := uuid.New()
	svc := &mockShowcaseService{
		getAvailableReposFn: func(_ context.Context, id uuid.UUID) ([]github.Repository, error) {
			assert.Equal(t, userID, id)
			return []github.Repository{{FullName: "user/repo"}}, nil
		},
	}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/repos/available", nil)
	req = injectClaims(req, userID)
	w := httptest.NewRecorder()

	h.HandleGetAvailableRepos(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleGetAvailableRepos_Unauthorized(t *testing.T) {
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/repos/available", nil)
	w := httptest.NewRecorder()

	h.HandleGetAvailableRepos(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- HandleSetShowcase ---

func TestHandleSetShowcase_Success(t *testing.T) {
	userID := uuid.New()
	svc := &mockShowcaseService{
		setShowcaseFn: func(_ context.Context, _ uuid.UUID, selections []domain.ShowcaseSelection) error {
			assert.Len(t, selections, 1)
			return nil
		},
	}
	h := NewShowcaseHandler(svc)

	body := `{"selections":[{"github_repo_id":123,"repo_full_name":"user/repo","tag":"coursework","description":"test"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/showcase", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	w := httptest.NewRecorder()

	h.HandleSetShowcase(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleSetShowcase_Unauthorized(t *testing.T) {
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/showcase", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	h.HandleSetShowcase(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleSetShowcase_InvalidJSON(t *testing.T) {
	userID := uuid.New()
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/showcase", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	w := httptest.NewRecorder()

	h.HandleSetShowcase(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleSetShowcase_ValidationFailsEmptySelections(t *testing.T) {
	userID := uuid.New()
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	// selections field is missing/null → validation should fail
	body := `{"selections":null}`
	req := httptest.NewRequest(http.MethodPost, "/api/showcase", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	w := httptest.NewRecorder()

	h.HandleSetShowcase(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- HandleGetShowcase ---

func TestHandleGetShowcase_Success(t *testing.T) {
	userID := uuid.New()
	svc := &mockShowcaseService{
		getShowcaseFn: func(_ context.Context, id uuid.UUID) ([]domain.ShowcaseRepo, error) {
			return []domain.ShowcaseRepo{}, nil
		},
	}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/showcase", nil)
	req = injectClaims(req, userID)
	w := httptest.NewRecorder()

	h.HandleGetShowcase(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleGetShowcase_Unauthorized(t *testing.T) {
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/showcase", nil)
	w := httptest.NewRecorder()

	h.HandleGetShowcase(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- HandleRemoveFromShowcase ---

func TestHandleRemoveFromShowcase_Success(t *testing.T) {
	userID := uuid.New()
	repoID := uuid.New()
	svc := &mockShowcaseService{
		removeFromShowcaseFn: func(_ context.Context, uid, rid uuid.UUID) error {
			assert.Equal(t, userID, uid)
			assert.Equal(t, repoID, rid)
			return nil
		},
	}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/showcase/"+repoID.String(), nil)
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleRemoveFromShowcase(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleRemoveFromShowcase_InvalidID(t *testing.T) {
	userID := uuid.New()
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/showcase/not-a-uuid", nil)
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", "not-a-uuid")
	w := httptest.NewRecorder()

	h.HandleRemoveFromShowcase(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleRemoveFromShowcase_NotFound(t *testing.T) {
	userID := uuid.New()
	repoID := uuid.New()
	svc := &mockShowcaseService{
		removeFromShowcaseFn: func(_ context.Context, _, _ uuid.UUID) error {
			return domain.ErrNotFound
		},
	}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/showcase/"+repoID.String(), nil)
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleRemoveFromShowcase(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --- HandleGetPublicRepo ---

func TestHandleGetPublicRepo_Success(t *testing.T) {
	repoID := uuid.New()
	svc := &mockShowcaseService{
		getPublicRepoFn: func(_ context.Context, id uuid.UUID) (*domain.ShowcaseRepo, error) {
			return &domain.ShowcaseRepo{ID: repoID, RepoFullName: "user/repo"}, nil
		},
	}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/repos/"+repoID.String(), nil)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleGetPublicRepo(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
}

func TestHandleGetPublicRepo_InvalidID(t *testing.T) {
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/repos/bad", nil)
	req = injectChiParam(req, "id", "bad")
	w := httptest.NewRecorder()

	h.HandleGetPublicRepo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleGetPublicRepo_NotFound(t *testing.T) {
	repoID := uuid.New()
	svc := &mockShowcaseService{
		getPublicRepoFn: func(_ context.Context, _ uuid.UUID) (*domain.ShowcaseRepo, error) {
			return nil, domain.ErrNotFound
		},
	}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/repos/%s", repoID), nil)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleGetPublicRepo(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --- HandleUpdateShowcaseRepo ---

func TestHandleUpdateShowcaseRepo_Success(t *testing.T) {
	userID := uuid.New()
	repoID := uuid.New()
	svc := &mockShowcaseService{
		updateRepoDescriptionFn: func(_ context.Context, uid, rid uuid.UUID, desc string) error {
			assert.Equal(t, userID, uid)
			assert.Equal(t, repoID, rid)
			assert.Equal(t, "My updated description", desc)
			return nil
		},
	}
	h := NewShowcaseHandler(svc)

	body := `{"description":"My updated description"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/showcase/"+repoID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleUpdateShowcaseRepo(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleUpdateShowcaseRepo_InvalidID(t *testing.T) {
	userID := uuid.New()
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/showcase/bad", strings.NewReader(`{}`))
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", "bad")
	w := httptest.NewRecorder()

	h.HandleUpdateShowcaseRepo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleUpdateShowcaseRepo_Unauthorized(t *testing.T) {
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/showcase/some-id", strings.NewReader(`{}`))
	req = injectChiParam(req, "id", uuid.New().String())
	w := httptest.NewRecorder()

	h.HandleUpdateShowcaseRepo(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleUpdateShowcaseRepo_InvalidJSON(t *testing.T) {
	userID := uuid.New()
	repoID := uuid.New()
	svc := &mockShowcaseService{}
	h := NewShowcaseHandler(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/showcase/"+repoID.String(), strings.NewReader("not json"))
	req = injectClaims(req, userID)
	req = injectChiParam(req, "id", repoID.String())
	w := httptest.NewRecorder()

	h.HandleUpdateShowcaseRepo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
