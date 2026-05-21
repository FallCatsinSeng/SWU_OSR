package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of GitHub event.
type EventType string

const (
	EventPush    EventType = "push"
	EventPR      EventType = "pull_request"
	EventRelease EventType = "release"
)

// ActivityLog stores a single GitHub webhook event.
type ActivityLog struct {
	ID             uuid.UUID       `json:"id"`
	UserID         uuid.UUID       `json:"user_id"`
	ShowcaseRepoID *uuid.UUID      `json:"showcase_repo_id"`
	RepoFullName   string          `json:"repo_full_name"`
	EventType      EventType       `json:"event_type"`
	Summary        string          `json:"summary"`
	Metadata       json.RawMessage `json:"metadata"`
	GitHubEventID  string          `json:"github_event_id"`
	CreatedAt      time.Time       `json:"created_at"`
}

// FeedParams holds pagination parameters for activity feeds.
type FeedParams struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

// FeedResult is the paginated activity feed response.
type FeedResult struct {
	Items      []ActivityItem `json:"items"`
	NextCursor string         `json:"next_cursor"`
	HasMore    bool           `json:"has_more"`
}

// ActivityItem is a single item in the activity feed.
type ActivityItem struct {
	ID           uuid.UUID       `json:"id"`
	UserID       uuid.UUID       `json:"user_id"`
	UserAlias    string          `json:"user_alias"`
	AvatarURL    string          `json:"avatar_url"`
	EventType    EventType       `json:"event_type"`
	RepoID       *uuid.UUID      `json:"repo_id"`
	RepoName     string          `json:"repo_name"`
	RepoFullName string          `json:"repo_full_name"`
	Summary      string          `json:"summary"`
	Metadata     json.RawMessage `json:"metadata"`
	CreatedAt    time.Time       `json:"created_at"`
}
