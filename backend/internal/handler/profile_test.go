package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mockProfileService implements service.ProfileService for testing.
type mockProfileService struct {
	getPublicProfileFn  func(ctx context.Context, alias string) (*service.PublicProfile, error)
	updateProfileFn     func(ctx context.Context, userID uuid.UUID, input service.UpdateProfileInput) error
	getRealIdentityFn   func(ctx context.Context, requesterID uuid.UUID, alias string) (*service.AcademicIdentity, error)
	getUserStatsFn      func(ctx context.Context, userID uuid.UUID) (*service.UserStats, error)
	listMembersFn       func(ctx context.Context) ([]*service.PublicProfile, error)
}

func (m *mockProfileService) GetPublicProfile(ctx context.Context, alias string) (*service.PublicProfile, error) {
	return m.getPublicProfileFn(ctx, alias)
}
func (m *mockProfileService) UpdateProfile(ctx context.Context, userID uuid.UUID, input service.UpdateProfileInput) error {
	return m.updateProfileFn(ctx, userID, input)
}
func (m *mockProfileService) GetRealIdentity(ctx context.Context, requesterID uuid.UUID, alias string) (*service.AcademicIdentity, error) {
	return m.getRealIdentityFn(ctx, requesterID, alias)
}
func (m *mockProfileService) GetUserStats(ctx context.Context, userID uuid.UUID) (*service.UserStats, error) {
	return m.getUserStatsFn(ctx, userID)
}
func (m *mockProfileService) ListMembers(ctx context.Context) ([]*service.PublicProfile, error) {
	return m.listMembersFn(ctx)
}

// --- HandleGetPublicProfile ---

func TestHandleGetPublicProfile_Success(t *testing.T) {
	svc := &mockProfileService{
		getPublicProfileFn: func(_ context.Context, alias string) (*service.PublicProfile, error) {
			assert.Equal(t, "johndoe", alias)
			return &service.PublicProfile{
				Alias: alias,
				Bio:   "A developer",
			}, nil
		},
	}
	h := NewProfileHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/johndoe", nil)
	req = injectChiParam(req, "alias", "johndoe")
	w := httptest.NewRecorder()

	h.HandleGetPublicProfile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleGetPublicProfile_NotFound(t *testing.T) {
	svc := &mockProfileService{
		getPublicProfileFn: func(_ context.Context, _ string) (*service.PublicProfile, error) {
			return nil, domain.ErrNotFound
		},
	}
	h := NewProfileHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/ghost", nil)
	req = injectChiParam(req, "alias", "ghost")
	w := httptest.NewRecorder()

	h.HandleGetPublicProfile(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --- HandleUpdateProfile ---

func TestHandleUpdateProfile_Success(t *testing.T) {
	userID := uuid.New()
	svc := &mockProfileService{
		updateProfileFn: func(_ context.Context, uid uuid.UUID, input service.UpdateProfileInput) error {
			assert.Equal(t, userID, uid)
			assert.Equal(t, "newbio", input.Bio)
			return nil
		},
	}
	h := NewProfileHandler(svc, nil)

	body := `{"bio":"newbio"}`
	req := httptest.NewRequest(http.MethodPut, "/api/profile", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	w := httptest.NewRecorder()

	h.HandleUpdateProfile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleUpdateProfile_Unauthorized(t *testing.T) {
	svc := &mockProfileService{}
	h := NewProfileHandler(svc, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/profile", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	h.HandleUpdateProfile(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleUpdateProfile_InvalidJSON(t *testing.T) {
	userID := uuid.New()
	svc := &mockProfileService{}
	h := NewProfileHandler(svc, nil)

	req := httptest.NewRequest(http.MethodPut, "/api/profile", strings.NewReader("invalid"))
	req = injectClaims(req, userID)
	w := httptest.NewRecorder()

	h.HandleUpdateProfile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleUpdateProfile_DuplicateAlias(t *testing.T) {
	userID := uuid.New()
	svc := &mockProfileService{
		updateProfileFn: func(_ context.Context, _ uuid.UUID, _ service.UpdateProfileInput) error {
			return domain.ErrDuplicateAlias
		},
	}
	h := NewProfileHandler(svc, nil)

	body := `{"alias":"taken_alias"}`
	req := httptest.NewRequest(http.MethodPut, "/api/profile", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = injectClaims(req, userID)
	w := httptest.NewRecorder()

	h.HandleUpdateProfile(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

// --- HandleListMembers ---

func TestHandleListMembers_Success(t *testing.T) {
	svc := &mockProfileService{
		listMembersFn: func(_ context.Context) ([]*service.PublicProfile, error) {
			return []*service.PublicProfile{
				{Alias: "alice"},
				{Alias: "bob"},
			}, nil
		},
	}
	h := NewProfileHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/members", nil)
	w := httptest.NewRecorder()

	h.HandleListMembers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleListMembers_Empty(t *testing.T) {
	svc := &mockProfileService{
		listMembersFn: func(_ context.Context) ([]*service.PublicProfile, error) {
			return []*service.PublicProfile{}, nil
		},
	}
	h := NewProfileHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/members", nil)
	w := httptest.NewRecorder()

	h.HandleListMembers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- HandleGetRealIdentity ---

func TestHandleGetRealIdentity_Success(t *testing.T) {
	requesterID := uuid.New()
	svc := &mockProfileService{
		getRealIdentityFn: func(_ context.Context, rid uuid.UUID, alias string) (*service.AcademicIdentity, error) {
			assert.Equal(t, requesterID, rid)
			assert.Equal(t, "johndoe", alias)
			return &service.AcademicIdentity{
				FullName: "John Doe",
				NIM:      "12345",
				Major:    "Computer Science",
				Semester: 4,
			}, nil
		},
	}
	h := NewProfileHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/johndoe/identity", nil)
	req = injectClaims(req, requesterID)
	req = injectChiParam(req, "alias", "johndoe")
	w := httptest.NewRecorder()

	h.HandleGetRealIdentity(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleGetRealIdentity_Unauthorized(t *testing.T) {
	svc := &mockProfileService{}
	h := NewProfileHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/johndoe/identity", nil)
	req = injectChiParam(req, "alias", "johndoe")
	w := httptest.NewRecorder()

	h.HandleGetRealIdentity(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleGetRealIdentity_TargetNotFound(t *testing.T) {
	requesterID := uuid.New()
	svc := &mockProfileService{
		getRealIdentityFn: func(_ context.Context, _ uuid.UUID, _ string) (*service.AcademicIdentity, error) {
			return nil, domain.ErrNotFound
		},
	}
	h := NewProfileHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/ghost/identity", nil)
	req = injectClaims(req, requesterID)
	req = injectChiParam(req, "alias", "ghost")
	w := httptest.NewRecorder()

	h.HandleGetRealIdentity(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
