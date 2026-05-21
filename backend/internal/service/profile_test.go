package service

import (
	"context"
	"testing"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockShowcaseRepo is a mock implementation of domain.ShowcaseRepository.
type mockShowcaseRepo struct {
	repos []domain.ShowcaseRepo
}

func newMockShowcaseRepo() *mockShowcaseRepo {
	return &mockShowcaseRepo{repos: []domain.ShowcaseRepo{}}
}

func (m *mockShowcaseRepo) Create(_ context.Context, repo *domain.ShowcaseRepo) error {
	m.repos = append(m.repos, *repo)
	return nil
}

func (m *mockShowcaseRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.ShowcaseRepo, error) {
	for i := range m.repos {
		if m.repos[i].ID == id && m.repos[i].DeletedAt == nil {
			return &m.repos[i], nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockShowcaseRepo) GetByUserID(_ context.Context, userID uuid.UUID) ([]domain.ShowcaseRepo, error) {
	var result []domain.ShowcaseRepo
	for _, r := range m.repos {
		if r.UserID == userID && r.DeletedAt == nil {
			result = append(result, r)
		}
	}
	if result == nil {
		result = []domain.ShowcaseRepo{}
	}
	return result, nil
}

func (m *mockShowcaseRepo) GetByUserAndRepoFullName(_ context.Context, userID uuid.UUID, repoFullName string) (*domain.ShowcaseRepo, error) {
	for i := range m.repos {
		if m.repos[i].UserID == userID && m.repos[i].RepoFullName == repoFullName && m.repos[i].DeletedAt == nil {
			return &m.repos[i], nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockShowcaseRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	now := time.Now()
	for i := range m.repos {
		if m.repos[i].ID == id {
			m.repos[i].DeletedAt = &now
		}
	}
	return nil
}

func (m *mockShowcaseRepo) SoftDeleteByUser(_ context.Context, userID uuid.UUID, repoID uuid.UUID) error {
	now := time.Now()
	for i := range m.repos {
		if m.repos[i].UserID == userID && m.repos[i].ID == repoID {
			m.repos[i].DeletedAt = &now
		}
	}
	return nil
}

func (m *mockShowcaseRepo) GetByUserAndRepoFullNameIncludeDeleted(_ context.Context, userID uuid.UUID, repoFullName string) (*domain.ShowcaseRepo, error) {
	for i := range m.repos {
		if m.repos[i].UserID == userID && m.repos[i].RepoFullName == repoFullName {
			return &m.repos[i], nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockShowcaseRepo) GetByUserAndGitHubRepoIDIncludeDeleted(_ context.Context, userID uuid.UUID, githubRepoID int64) (*domain.ShowcaseRepo, error) {
	for i := range m.repos {
		if m.repos[i].UserID == userID && m.repos[i].GitHubRepoID == githubRepoID {
			return &m.repos[i], nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockShowcaseRepo) Restore(_ context.Context, repo *domain.ShowcaseRepo) error {
	for i := range m.repos {
		if m.repos[i].ID == repo.ID {
			m.repos[i] = *repo
			return nil
		}
	}
	return nil
}

func (m *mockShowcaseRepo) UpdateDescription(_ context.Context, repoID uuid.UUID, description string) error {
	for i := range m.repos {
		if m.repos[i].ID == repoID {
			m.repos[i].Description = description
		}
	}
	return nil
}

// mockActivityRepo is a mock implementation of domain.ActivityRepository.
type mockActivityRepo struct {
	logs []domain.ActivityLog
}

func newMockActivityRepo() *mockActivityRepo {
	return &mockActivityRepo{logs: []domain.ActivityLog{}}
}

func (m *mockActivityRepo) Insert(_ context.Context, log *domain.ActivityLog) error {
	for _, l := range m.logs {
		if l.GitHubEventID == log.GitHubEventID {
			return nil // idempotent
		}
	}
	m.logs = append(m.logs, *log)
	return nil
}

func (m *mockActivityRepo) GetFeed(_ context.Context, cursor time.Time, limit int) ([]domain.ActivityItem, error) {
	var items []domain.ActivityItem
	for _, l := range m.logs {
		if l.CreatedAt.Before(cursor) {
			items = append(items, domain.ActivityItem{
				ID:        l.ID,
				UserID:    l.UserID,
				EventType: l.EventType,
				Summary:   l.Summary,
				Metadata:  l.Metadata,
				CreatedAt: l.CreatedAt,
			})
		}
	}
	if len(items) > limit {
		items = items[:limit]
	}
	if items == nil {
		items = []domain.ActivityItem{}
	}
	return items, nil
}

func (m *mockActivityRepo) GetUserFeed(_ context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]domain.ActivityItem, error) {
	var items []domain.ActivityItem
	for _, l := range m.logs {
		if l.UserID == userID && l.CreatedAt.Before(cursor) {
			items = append(items, domain.ActivityItem{
				ID:        l.ID,
				UserID:    l.UserID,
				EventType: l.EventType,
				Summary:   l.Summary,
				Metadata:  l.Metadata,
				CreatedAt: l.CreatedAt,
			})
		}
	}
	if len(items) > limit {
		items = items[:limit]
	}
	if items == nil {
		items = []domain.ActivityItem{}
	}
	return items, nil
}

func (m *mockActivityRepo) GetRepoFeed(_ context.Context, showcaseRepoID uuid.UUID, cursor time.Time, limit int) ([]domain.ActivityItem, error) {
	var items []domain.ActivityItem
	for _, l := range m.logs {
		if l.ShowcaseRepoID != nil && *l.ShowcaseRepoID == showcaseRepoID && l.CreatedAt.Before(cursor) {
			items = append(items, domain.ActivityItem{
				ID:        l.ID,
				UserID:    l.UserID,
				EventType: l.EventType,
				Summary:   l.Summary,
				Metadata:  l.Metadata,
				CreatedAt: l.CreatedAt,
			})
		}
	}
	if len(items) > limit {
		items = items[:limit]
	}
	if items == nil {
		items = []domain.ActivityItem{}
	}
	return items, nil
}

func (m *mockActivityRepo) GetByGitHubEventID(_ context.Context, eventID string) (*domain.ActivityLog, error) {
	for i := range m.logs {
		if m.logs[i].GitHubEventID == eventID {
			return &m.logs[i], nil
		}
	}
	return nil, domain.ErrNotFound
}

func TestGetPublicProfile_NeverExposesPrivateData(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	activityRepo := newMockActivityRepo()

	userID := uuid.New()
	user := &domain.User{
		ID:             userID,
		NIM:            "2021001",
		FullName:       "John Doe",
		Major:          "Computer Science",
		Semester:       6,
		Alias:          "johndoe",
		Bio:            "A developer",
		AvatarURL:      "https://example.com/avatar.png",
		GitHubUsername: "johndoe",
		Role:           domain.RoleStudent,
		CreatedAt:      time.Now(),
	}
	userRepo.users["2021001"] = user

	svc := NewProfileService(userRepo, showcaseRepo, activityRepo, nil, nil)

	profile, err := svc.GetPublicProfile(context.Background(), "johndoe")
	require.NoError(t, err)

	// Public data should be present
	assert.Equal(t, "johndoe", profile.Alias)
	assert.Equal(t, "A developer", profile.Bio)
	assert.Equal(t, "https://example.com/avatar.png", profile.AvatarURL)

	// The PublicProfile type does not contain NIM, FullName, Major, or Semester fields
	assert.NotNil(t, profile.Stats)
	assert.NotNil(t, profile.ShowcaseRepos)
}

func TestGetRealIdentity_FacultySucceeds(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	activityRepo := newMockActivityRepo()

	facultyID := uuid.New()
	studentID := uuid.New()

	faculty := &domain.User{
		ID:    facultyID,
		NIM:   "FAC001",
		Alias: "prof_smith",
		Role:  domain.RoleFaculty,
	}
	student := &domain.User{
		ID:       studentID,
		NIM:      "2021001",
		FullName: "John Doe",
		Major:    "Computer Science",
		Semester: 6,
		Alias:    "johndoe",
		Role:     domain.RoleStudent,
	}

	userRepo.users["FAC001"] = faculty
	userRepo.users["2021001"] = student

	svc := NewProfileService(userRepo, showcaseRepo, activityRepo, nil, nil)

	identity, err := svc.GetRealIdentity(context.Background(), facultyID, "johndoe")
	require.NoError(t, err)
	assert.Equal(t, "John Doe", identity.FullName)
	assert.Equal(t, "2021001", identity.NIM)
	assert.Equal(t, "Computer Science", identity.Major)
	assert.Equal(t, 6, identity.Semester)
}

func TestGetRealIdentity_StudentSucceeds(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	activityRepo := newMockActivityRepo()

	studentID := uuid.New()
	otherStudentID := uuid.New()

	student := &domain.User{
		ID:    studentID,
		NIM:   "2021001",
		Alias: "student1",
		Role:  domain.RoleStudent,
	}
	otherStudent := &domain.User{
		ID:       otherStudentID,
		NIM:      "2021002",
		FullName: "Jane Doe",
		Major:    "Mathematics",
		Semester: 4,
		Alias:    "student2",
		Role:     domain.RoleStudent,
	}

	userRepo.users["2021001"] = student
	userRepo.users["2021002"] = otherStudent

	svc := NewProfileService(userRepo, showcaseRepo, activityRepo, nil, nil)

	identity, err := svc.GetRealIdentity(context.Background(), studentID, "student2")
	require.NoError(t, err)
	assert.Equal(t, "Jane Doe", identity.FullName)
	assert.Equal(t, "2021002", identity.NIM)
	assert.Equal(t, "Mathematics", identity.Major)
	assert.Equal(t, 4, identity.Semester)
}
