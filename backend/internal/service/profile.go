package service

import (
	"context"
	"regexp"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
)

// PublicProfile is the public-facing user profile. Never includes NIM, full_name, major, semester.
type PublicProfile struct {
	ID             uuid.UUID            `json:"id"`
	Alias          string               `json:"alias"`
	Bio            string               `json:"bio"`
	AvatarURL      string               `json:"avatar_url"`
	GitHubUsername string               `json:"github_username"`
	Role           domain.Role          `json:"role"`
	ShowcaseRepos  []domain.ShowcaseRepo `json:"showcase_repos"`
	Stats          *UserStats           `json:"stats"`
	CreatedAt      time.Time            `json:"created_at"`
}

// AcademicIdentity holds private academic information only visible to faculty/admin.
type AcademicIdentity struct {
	FullName string `json:"full_name"`
	NIM      string `json:"nim"`
	Major    string `json:"major"`
	Semester int    `json:"semester"`
}

// UserStats holds computed activity statistics for a user.
type UserStats struct {
	TotalCommits  int      `json:"total_commits"`
	TotalRepos    int      `json:"total_repos"`
	Languages     []string `json:"languages"`
	ActiveDays    int      `json:"active_days"`
	CurrentStreak int      `json:"current_streak"`
}

// ProfileService defines the profile service interface.
type ProfileService interface {
	GetPublicProfile(ctx context.Context, alias string) (*PublicProfile, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) error
	GetRealIdentity(ctx context.Context, requesterID uuid.UUID, alias string) (*AcademicIdentity, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
	ListMembers(ctx context.Context) ([]*PublicProfile, error)
}

// UpdateProfileInput holds the input for updating a profile.
type UpdateProfileInput struct {
	Alias     string `json:"alias"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

var aliasRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`)

// profileService is the concrete implementation.
type profileService struct {
	userRepo     domain.UserRepository
	showcaseRepo domain.ShowcaseRepository
	activityRepo domain.ActivityRepository
}

// NewProfileService creates a new profile service.
func NewProfileService(
	userRepo domain.UserRepository,
	showcaseRepo domain.ShowcaseRepository,
	activityRepo domain.ActivityRepository,
) ProfileService {
	return &profileService{
		userRepo:     userRepo,
		showcaseRepo: showcaseRepo,
		activityRepo: activityRepo,
	}
}

// GetPublicProfile returns the public profile for a user by alias.
// It NEVER includes NIM, full_name, major, or semester.
func (s *profileService) GetPublicProfile(ctx context.Context, alias string) (*PublicProfile, error) {
	user, err := s.userRepo.GetByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}

	repos, err := s.showcaseRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	stats, err := s.GetUserStats(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &PublicProfile{
		ID:             user.ID,
		Alias:          user.Alias,
		Bio:            user.Bio,
		AvatarURL:      user.AvatarURL,
		GitHubUsername: user.GitHubUsername,
		Role:           user.Role,
		ShowcaseRepos:  repos,
		Stats:          stats,
		CreatedAt:      user.CreatedAt,
	}, nil
}

// UpdateProfile validates and persists profile changes.
func (s *profileService) UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) error {
	if input.Alias != "" {
		if !aliasRegex.MatchString(input.Alias) {
			return domain.ErrDuplicateAlias
		}

		// Check uniqueness
		existing, err := s.userRepo.GetByAlias(ctx, input.Alias)
		if err == nil && existing.ID != userID {
			return domain.ErrDuplicateAlias
		}
		if err != nil && err != domain.ErrNotFound {
			return err
		}
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if input.Alias != "" {
		user.Alias = input.Alias
	}
	if input.Bio != "" {
		user.Bio = input.Bio
	}
	if input.AvatarURL != "" {
		user.AvatarURL = input.AvatarURL
	}
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

// GetRealIdentity returns academic identity for any authenticated user.
func (s *profileService) GetRealIdentity(ctx context.Context, requesterID uuid.UUID, alias string) (*AcademicIdentity, error) {
	_, err := s.userRepo.GetByID(ctx, requesterID)
	if err != nil {
		return nil, err
	}

	target, err := s.userRepo.GetByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}

	return &AcademicIdentity{
		FullName: target.FullName,
		NIM:      target.NIM,
		Major:    target.Major,
		Semester: target.Semester,
	}, nil
}

// GetUserStats computes activity statistics for a user.
func (s *profileService) GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error) {
	repos, err := s.showcaseRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Collect languages from showcase repos
	langSet := make(map[string]struct{})
	for _, r := range repos {
		if r.Language != "" {
			langSet[r.Language] = struct{}{}
		}
	}
	languages := make([]string, 0, len(langSet))
	for l := range langSet {
		languages = append(languages, l)
	}

	// Get activity feed to compute stats
	feed, err := s.activityRepo.GetUserFeed(ctx, userID, time.Now().Add(time.Second), 1000)
	if err != nil {
		return nil, err
	}

	totalCommits := 0
	daySet := make(map[string]struct{})
	for _, item := range feed {
		if item.EventType == domain.EventPush {
			totalCommits++
		}
		day := item.CreatedAt.Format("2006-01-02")
		daySet[day] = struct{}{}
	}

	// Calculate streak
	streak := 0
	today := time.Now().Truncate(24 * time.Hour)
	for i := 0; ; i++ {
		day := today.AddDate(0, 0, -i).Format("2006-01-02")
		if _, ok := daySet[day]; ok {
			streak++
		} else {
			break
		}
	}

	return &UserStats{
		TotalCommits:  totalCommits,
		TotalRepos:    len(repos),
		Languages:     languages,
		ActiveDays:    len(daySet),
		CurrentStreak: streak,
	}, nil
}


// ListMembers returns all active users as public profiles.
func (s *profileService) ListMembers(ctx context.Context) ([]*PublicProfile, error) {
	users, err := s.userRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	profiles := make([]*PublicProfile, 0, len(users))
	for _, user := range users {
		repos, err := s.showcaseRepo.GetByUserID(ctx, user.ID)
		if err != nil {
			repos = nil
		}

		stats, err := s.GetUserStats(ctx, user.ID)
		if err != nil {
			stats = nil
		}

		profiles = append(profiles, &PublicProfile{
			ID:             user.ID,
			Alias:          user.Alias,
			Bio:            user.Bio,
			AvatarURL:      user.AvatarURL,
			GitHubUsername: user.GitHubUsername,
			Role:           user.Role,
			ShowcaseRepos:  repos,
			Stats:          stats,
			CreatedAt:      user.CreatedAt,
		})
	}

	return profiles, nil
}
