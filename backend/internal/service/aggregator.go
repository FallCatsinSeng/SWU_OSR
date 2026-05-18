package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
)

// AggregatorService defines the aggregator service interface.
type AggregatorService interface {
	ProcessWebhook(ctx context.Context, payload []byte, signature, eventType, deliveryID string) error
	GetActivityFeed(ctx context.Context, params domain.FeedParams) (*domain.FeedResult, error)
	GetUserActivity(ctx context.Context, userID uuid.UUID, params domain.FeedParams) (*domain.FeedResult, error)
}

// aggregatorService is the concrete implementation.
type aggregatorService struct {
	activityRepo  domain.ActivityRepository
	userRepo      domain.UserRepository
	showcaseRepo  domain.ShowcaseRepository
	webhookSecret []byte
}

// NewAggregatorService creates a new aggregator service.
func NewAggregatorService(
	activityRepo domain.ActivityRepository,
	userRepo domain.UserRepository,
	showcaseRepo domain.ShowcaseRepository,
	webhookSecret string,
) AggregatorService {
	return &aggregatorService{
		activityRepo:  activityRepo,
		userRepo:      userRepo,
		showcaseRepo:  showcaseRepo,
		webhookSecret: []byte(webhookSecret),
	}
}

// ProcessWebhook verifies the HMAC signature, parses the event, and inserts an activity log.
func (s *aggregatorService) ProcessWebhook(ctx context.Context, payload []byte, signature, eventType, deliveryID string) error {
	// Verify HMAC-SHA256 signature
	if !s.verifySignature(payload, signature) {
		return domain.ErrInvalidSignature
	}

	// Parse the event
	var repoFullName, username, summary string
	var metadata json.RawMessage

	switch domain.EventType(eventType) {
	case domain.EventPush:
		repoFullName, username, summary, metadata = s.parsePushEvent(payload)
	case domain.EventPR:
		repoFullName, username, summary, metadata = s.parsePREvent(payload)
	case domain.EventRelease:
		repoFullName, username, summary, metadata = s.parseReleaseEvent(payload)
	default:
		// Unsupported event type, ignore
		return nil
	}

	if username == "" || repoFullName == "" {
		return nil
	}

	// Resolve user by GitHub username
	user, err := s.userRepo.GetByGitHubUsername(ctx, username)
	if err != nil {
		// User not found, ignore
		return nil
	}

	// Resolve showcase repo
	showcaseRepo, err := s.showcaseRepo.GetByUserAndRepoFullName(ctx, user.ID, repoFullName)
	if err != nil {
		// Not in showcase, ignore
		return nil
	}

	// Check duplicate via github_event_id
	if deliveryID != "" {
		existing, err := s.activityRepo.GetByGitHubEventID(ctx, deliveryID)
		if err == nil && existing != nil {
			// Already processed, skip
			return nil
		}
	}

	// Insert activity log
	log := &domain.ActivityLog{
		ID:             uuid.New(),
		UserID:         user.ID,
		ShowcaseRepoID: showcaseRepo.ID,
		EventType:      domain.EventType(eventType),
		Summary:        summary,
		Metadata:       metadata,
		GitHubEventID:  deliveryID,
		CreatedAt:      time.Now(),
	}

	return s.activityRepo.Insert(ctx, log)
}

// GetActivityFeed returns a paginated activity feed.
func (s *aggregatorService) GetActivityFeed(ctx context.Context, params domain.FeedParams) (*domain.FeedResult, error) {
	limit := clampLimit(params.Limit)
	cursor := decodeCursor(params.Cursor)

	items, err := s.activityRepo.GetFeed(ctx, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}

	var nextCursor string
	if hasMore && len(items) > 0 {
		nextCursor = encodeCursor(items[len(items)-1].CreatedAt)
	}

	return &domain.FeedResult{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// GetUserActivity returns a paginated activity feed filtered by user.
func (s *aggregatorService) GetUserActivity(ctx context.Context, userID uuid.UUID, params domain.FeedParams) (*domain.FeedResult, error) {
	limit := clampLimit(params.Limit)
	cursor := decodeCursor(params.Cursor)

	items, err := s.activityRepo.GetUserFeed(ctx, userID, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}

	var nextCursor string
	if hasMore && len(items) > 0 {
		nextCursor = encodeCursor(items[len(items)-1].CreatedAt)
	}

	return &domain.FeedResult{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// verifySignature checks the HMAC-SHA256 signature.
func (s *aggregatorService) verifySignature(payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, s.webhookSecret)
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// parsePushEvent extracts data from a push event payload.
func (s *aggregatorService) parsePushEvent(payload []byte) (repoFullName, username, summary string, metadata json.RawMessage) {
	var event struct {
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
		Pusher struct {
			Name string `json:"name"`
		} `json:"pusher"`
		Ref     string `json:"ref"`
		Commits []struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"commits"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return "", "", "", nil
	}

	commitCount := len(event.Commits)
	branch := event.Ref
	if strings.HasPrefix(branch, "refs/heads/") {
		branch = strings.TrimPrefix(branch, "refs/heads/")
	}

	summary = fmt.Sprintf("Pushed %d commit(s) to %s", commitCount, branch)

	// Build metadata with up to 5 commits
	commits := event.Commits
	if len(commits) > 5 {
		commits = commits[:5]
	}
	meta := map[string]interface{}{
		"ref":          event.Ref,
		"commit_count": commitCount,
		"commits":      commits,
	}
	metadata, _ = json.Marshal(meta)

	return event.Repository.FullName, event.Pusher.Name, summary, metadata
}

// parsePREvent extracts data from a pull_request event payload.
func (s *aggregatorService) parsePREvent(payload []byte) (repoFullName, username, summary string, metadata json.RawMessage) {
	var event struct {
		Action string `json:"action"`
		Number int    `json:"number"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
		PullRequest struct {
			Title string `json:"title"`
			User  struct {
				Login string `json:"login"`
			} `json:"user"`
		} `json:"pull_request"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return "", "", "", nil
	}

	summary = fmt.Sprintf("PR #%d %s: %s", event.Number, event.Action, event.PullRequest.Title)

	meta := map[string]interface{}{
		"action": event.Action,
		"number": event.Number,
		"title":  event.PullRequest.Title,
	}
	metadata, _ = json.Marshal(meta)

	return event.Repository.FullName, event.PullRequest.User.Login, summary, metadata
}

// parseReleaseEvent extracts data from a release event payload.
func (s *aggregatorService) parseReleaseEvent(payload []byte) (repoFullName, username, summary string, metadata json.RawMessage) {
	var event struct {
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
		Release struct {
			TagName string `json:"tag_name"`
			Name    string `json:"name"`
			Author  struct {
				Login string `json:"login"`
			} `json:"author"`
		} `json:"release"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return "", "", "", nil
	}

	summary = fmt.Sprintf("Released %s (%s)", event.Release.Name, event.Release.TagName)

	meta := map[string]interface{}{
		"tag_name": event.Release.TagName,
		"name":     event.Release.Name,
	}
	metadata, _ = json.Marshal(meta)

	return event.Repository.FullName, event.Release.Author.Login, summary, metadata
}

// clampLimit ensures the limit is between 1 and 50, defaulting to 20.
func clampLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 50 {
		return 50
	}
	return limit
}

// decodeCursor decodes a base64-encoded RFC3339 timestamp cursor.
func decodeCursor(cursor string) time.Time {
	if cursor == "" {
		return time.Now().Add(time.Second)
	}
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return time.Now().Add(time.Second)
	}
	t, err := time.Parse(time.RFC3339Nano, string(data))
	if err != nil {
		return time.Now().Add(time.Second)
	}
	return t
}

// encodeCursor encodes a time as a base64 RFC3339Nano string.
func encodeCursor(t time.Time) string {
	return base64.StdEncoding.EncodeToString([]byte(t.Format(time.RFC3339Nano)))
}
