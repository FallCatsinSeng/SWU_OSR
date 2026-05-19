package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/github"
	"github.com/google/uuid"
)

// ShowcaseService defines the showcase service interface.
type ShowcaseService interface {
	GetAvailableRepos(ctx context.Context, userID uuid.UUID) ([]github.Repository, error)
	SetShowcase(ctx context.Context, userID uuid.UUID, selections []domain.ShowcaseSelection) error
	GetShowcase(ctx context.Context, userID uuid.UUID) ([]domain.ShowcaseRepo, error)
	UpdateShowcase(ctx context.Context, userID uuid.UUID, selections []domain.ShowcaseSelection) error
	RemoveFromShowcase(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) error
	UpdateRepoDescription(ctx context.Context, userID uuid.UUID, repoID uuid.UUID, description string) error
	SetAggregatorService(agg AggregatorService)
}

// showcaseService is the concrete implementation.
type showcaseService struct {
	showcaseRepo   domain.ShowcaseRepository
	userRepo       domain.UserRepository
	githubSvc      github.Service
	encryptionKey  []byte
	webhookURL     string
	webhookSecret  string
	aggregatorSvc  AggregatorService
}

// NewShowcaseService creates a new showcase service.
func NewShowcaseService(
	showcaseRepo domain.ShowcaseRepository,
	userRepo domain.UserRepository,
	githubSvc github.Service,
	encryptionKey []byte,
	webhookURL string,
	webhookSecret string,
) ShowcaseService {
	return &showcaseService{
		showcaseRepo:  showcaseRepo,
		userRepo:      userRepo,
		githubSvc:     githubSvc,
		encryptionKey: encryptionKey,
		webhookURL:    webhookURL,
		webhookSecret: webhookSecret,
	}
}

// SetAggregatorService sets the aggregator service for auto-sync on showcase add.
// This is set after construction to avoid circular initialization.
func (s *showcaseService) SetAggregatorService(agg AggregatorService) {
	s.aggregatorSvc = agg
}

// validAcademicTags is the set of valid academic tags.
var validAcademicTags = map[domain.AcademicTag]bool{
	domain.TagCoursework:       true,
	domain.TagThesis:           true,
	domain.TagHackathon:        true,
	domain.TagPersonalResearch: true,
	domain.TagTeamProject:      true,
}

// GetAvailableRepos decrypts the user's GitHub token and lists their repos.
func (s *showcaseService) GetAvailableRepos(ctx context.Context, userID uuid.UUID) ([]github.Repository, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	token, err := Decrypt(user.GitHubToken, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("decrypting github token: %w", err)
	}

	return s.githubSvc.ListRepos(ctx, token)
}

// SetShowcase validates selections and adds them to the user's showcase (appending, not replacing).
func (s *showcaseService) SetShowcase(ctx context.Context, userID uuid.UUID, selections []domain.ShowcaseSelection) error {
	if len(selections) > 20 {
		return fmt.Errorf("maximum 20 showcase repos allowed")
	}

	for _, sel := range selections {
		if !validAcademicTags[sel.Tag] {
			return fmt.Errorf("invalid academic tag: %s", sel.Tag)
		}
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	token, err := Decrypt(user.GitHubToken, s.encryptionKey)
	if err != nil {
		return fmt.Errorf("decrypting github token: %w", err)
	}

	// Get existing repos to check total count and avoid duplicates
	existing, err := s.showcaseRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if len(existing)+len(selections) > 20 {
		return fmt.Errorf("maximum 20 showcase repos allowed (currently %d)", len(existing))
	}

	// Build set of already-showcased full names and repo IDs to skip duplicates
	existingFullNames := make(map[string]bool)
	existingRepoIDs := make(map[int64]bool)
	for _, repo := range existing {
		existingFullNames[repo.RepoFullName] = true
		existingRepoIDs[repo.GitHubRepoID] = true
	}

	// Fetch all repos from GitHub to get metadata (html_url, description, language)
	ghRepos, _ := s.githubSvc.ListRepos(ctx, token)
	repoMetaMap := make(map[string]github.Repository)
	for _, r := range ghRepos {
		repoMetaMap[r.FullName] = r
	}

	for _, sel := range selections {
		// Skip if already in showcase (active) — check both full_name and repo_id
		if existingFullNames[sel.FullName] || existingRepoIDs[sel.RepoID] {
			continue
		}

		// Check if this repo was previously soft-deleted — if so, restore it
		// We check by both full_name and github_repo_id (the DB constraint is on github_repo_id)
		existingDeleted, _ := s.showcaseRepo.GetByUserAndRepoFullNameIncludeDeleted(ctx, userID, sel.FullName)
		if existingDeleted == nil {
			// Also try by github_repo_id directly
			existingDeleted, _ = s.showcaseRepo.GetByUserAndGitHubRepoIDIncludeDeleted(ctx, userID, sel.RepoID)
		}
		if existingDeleted != nil {
			// Restore the soft-deleted entry
			existingDeleted.DeletedAt = nil
			existingDeleted.AcademicTag = sel.Tag
			existingDeleted.RepoName = sel.RepoName
			existingDeleted.RepoFullName = sel.FullName
			existingDeleted.UpdatedAt = time.Now()

			// Update metadata
			if meta, ok := repoMetaMap[sel.FullName]; ok {
				existingDeleted.Description = meta.Description
				existingDeleted.Language = meta.Language
				existingDeleted.HTMLURL = meta.HTMLURL
			}
			if existingDeleted.HTMLURL == "" {
				existingDeleted.HTMLURL = fmt.Sprintf("https://github.com/%s", sel.FullName)
			}

			// Re-register webhook
			parts := strings.SplitN(sel.FullName, "/", 2)
			if len(parts) == 2 {
				id, whErr := s.githubSvc.RegisterWebhook(ctx, token, parts[0], parts[1], s.webhookURL, s.webhookSecret)
				if whErr == nil {
					existingDeleted.WebhookID = &id
				}
			}

			if err := s.showcaseRepo.Restore(ctx, existingDeleted); err != nil {
				continue
			}

			// Auto-sync activity for restored repo (background, non-blocking)
			if s.aggregatorSvc != nil {
				go func(uid uuid.UUID, repoID uuid.UUID) {
					syncCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
					defer cancel()
					_, _ = s.aggregatorSvc.SyncRepoActivity(syncCtx, uid, repoID)
				}(userID, existingDeleted.ID)
			}
			continue
		}

		// New repo — register webhook and create
		parts := strings.SplitN(sel.FullName, "/", 2)
		var webhookID *int64
		if len(parts) == 2 {
			id, whErr := s.githubSvc.RegisterWebhook(ctx, token, parts[0], parts[1], s.webhookURL, s.webhookSecret)
			if whErr == nil {
				webhookID = &id
			}
		}

		// Get metadata from GitHub API response
		var description, language, htmlURL string
		if meta, ok := repoMetaMap[sel.FullName]; ok {
			description = meta.Description
			language = meta.Language
			htmlURL = meta.HTMLURL
		}
		if htmlURL == "" {
			htmlURL = fmt.Sprintf("https://github.com/%s", sel.FullName)
		}

		repo := &domain.ShowcaseRepo{
			ID:           uuid.New(),
			UserID:       userID,
			GitHubRepoID: sel.RepoID,
			RepoName:     sel.RepoName,
			RepoFullName: sel.FullName,
			Description:  description,
			Language:     language,
			HTMLURL:      htmlURL,
			AcademicTag:  sel.Tag,
			WebhookID:    webhookID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := s.showcaseRepo.Create(ctx, repo); err != nil {
			// If it fails (e.g. constraint), skip this repo but don't abort
			continue
		}

		// Auto-sync activity for newly added repo (background, non-blocking)
		if s.aggregatorSvc != nil {
			go func(uid uuid.UUID, repoID uuid.UUID) {
				syncCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()
				_, _ = s.aggregatorSvc.SyncRepoActivity(syncCtx, uid, repoID)
			}(userID, repo.ID)
		}
	}

	return nil
}

// GetShowcase returns the user's showcase repos from the database.
func (s *showcaseService) GetShowcase(ctx context.Context, userID uuid.UUID) ([]domain.ShowcaseRepo, error) {
	return s.showcaseRepo.GetByUserID(ctx, userID)
}

// UpdateShowcase is the same as SetShowcase (replace all).
func (s *showcaseService) UpdateShowcase(ctx context.Context, userID uuid.UUID, selections []domain.ShowcaseSelection) error {
	return s.SetShowcase(ctx, userID, selections)
}

// RemoveFromShowcase soft-deletes a specific repo and removes its webhook.
func (s *showcaseService) RemoveFromShowcase(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Find the repo in user's showcase
	repos, err := s.showcaseRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var target *domain.ShowcaseRepo
	for i := range repos {
		if repos[i].ID == repoID {
			target = &repos[i]
			break
		}
	}
	if target == nil {
		return domain.ErrNotFound
	}

	// Remove webhook if it exists
	if target.WebhookID != nil {
		token, err := Decrypt(user.GitHubToken, s.encryptionKey)
		if err == nil {
			parts := strings.SplitN(target.RepoFullName, "/", 2)
			if len(parts) == 2 {
				_ = s.githubSvc.RemoveWebhook(ctx, token, parts[0], parts[1], *target.WebhookID)
			}
		}
	}

	return s.showcaseRepo.SoftDeleteByUser(ctx, userID, repoID)
}


// UpdateRepoDescription updates the description of a specific showcase repo.
func (s *showcaseService) UpdateRepoDescription(ctx context.Context, userID uuid.UUID, repoID uuid.UUID, description string) error {
	repos, err := s.showcaseRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var target *domain.ShowcaseRepo
	for i := range repos {
		if repos[i].ID == repoID {
			target = &repos[i]
			break
		}
	}
	if target == nil {
		return domain.ErrNotFound
	}

	target.Description = description
	target.UpdatedAt = time.Now()
	return s.showcaseRepo.UpdateDescription(ctx, repoID, description)
}
