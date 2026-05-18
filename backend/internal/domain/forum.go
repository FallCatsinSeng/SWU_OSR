package domain

import (
	"time"

	"github.com/google/uuid"
)

// Thread represents a discussion thread on a showcase repo.
type Thread struct {
	ID             uuid.UUID  `json:"id"`
	ShowcaseRepoID uuid.UUID  `json:"showcase_repo_id"`
	AuthorID       uuid.UUID  `json:"author_id"`
	Title          string     `json:"title"`
	Body           string     `json:"body"`
	CommentCount   int        `json:"comment_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"-"`
}

// Comment represents a reply in a thread.
type Comment struct {
	ID        uuid.UUID  `json:"id"`
	ThreadID  uuid.UUID  `json:"thread_id"`
	AuthorID  uuid.UUID  `json:"author_id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// Notification represents an internal notification for a user.
type Notification struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Type        string    `json:"type"`
	ReferenceID uuid.UUID `json:"reference_id"`
	Message     string    `json:"message"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateThreadInput is the input for creating a new thread.
type CreateThreadInput struct {
	RepoID   uuid.UUID `json:"repo_id" validate:"required"`
	AuthorID uuid.UUID `json:"author_id" validate:"required"`
	Title    string    `json:"title" validate:"required,min=5,max=255"`
	Body     string    `json:"body" validate:"required,min=1,max=10000"`
}

// CreateCommentInput is the input for creating a new comment.
type CreateCommentInput struct {
	ThreadID uuid.UUID  `json:"thread_id" validate:"required"`
	AuthorID uuid.UUID  `json:"author_id" validate:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
	Body     string     `json:"body" validate:"required,min=1,max=10000"`
}

// PaginationParams holds pagination parameters for thread listing.
type PaginationParams struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

// ThreadList is the paginated thread list response.
type ThreadList struct {
	Threads    []Thread `json:"threads"`
	NextCursor string   `json:"next_cursor"`
	HasMore    bool     `json:"has_more"`
}
