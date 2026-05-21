package service

import (
	"context"
	"fmt"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/sanitize"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ForumService defines the forum service interface.
type ForumService interface {
	ListThreads(ctx context.Context, repoID uuid.UUID, params domain.PaginationParams) (*domain.ThreadList, error)
	CreateThread(ctx context.Context, input domain.CreateThreadInput) (*domain.Thread, error)
	GetThread(ctx context.Context, threadID uuid.UUID) (*domain.Thread, []domain.Comment, error)
	CreateComment(ctx context.Context, input domain.CreateCommentInput) (*domain.Comment, error)
	ListNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error)
	MarkNotificationRead(ctx context.Context, userID uuid.UUID, notifID uuid.UUID) error
}

// forumService is the concrete implementation.
type forumService struct {
	threadRepo   domain.ThreadRepository
	commentRepo  domain.CommentRepository
	notifRepo    domain.NotificationRepository
	showcaseRepo domain.ShowcaseRepository
	userRepo     domain.UserRepository
	logger       *zap.Logger
}

// NewForumService creates a new forum service.
func NewForumService(
	threadRepo domain.ThreadRepository,
	commentRepo domain.CommentRepository,
	notifRepo domain.NotificationRepository,
	showcaseRepo domain.ShowcaseRepository,
	userRepo domain.UserRepository,
	logger *zap.Logger,
) ForumService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &forumService{
		threadRepo:   threadRepo,
		commentRepo:  commentRepo,
		notifRepo:    notifRepo,
		showcaseRepo: showcaseRepo,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// ListThreads returns paginated threads for a showcase repo.
func (s *forumService) ListThreads(ctx context.Context, repoID uuid.UUID, params domain.PaginationParams) (*domain.ThreadList, error) {
	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	cursor := decodeCursor(params.Cursor)

	threads, err := s.threadRepo.ListByRepoID(ctx, repoID, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	hasMore := len(threads) > limit
	if hasMore {
		threads = threads[:limit]
	}

	var nextCursor string
	if hasMore && len(threads) > 0 {
		nextCursor = encodeCursor(threads[len(threads)-1].CreatedAt)
	}

	return &domain.ThreadList{
		Threads:    threads,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// CreateThread validates input, creates a thread, and notifies the repo owner.
func (s *forumService) CreateThread(ctx context.Context, input domain.CreateThreadInput) (*domain.Thread, error) {
	// Sanitize: strip HTML tags from user input to prevent stored XSS
	input.Title = sanitize.StripHTML(input.Title)
	input.Body = sanitize.StripHTMLPreserveWhitespace(input.Body)

	// Validate title length
	if len(input.Title) < 5 {
		return nil, fmt.Errorf("title must be at least 5 characters")
	}
	if len(input.Title) > 255 {
		return nil, fmt.Errorf("title must not exceed 255 characters")
	}

	// Validate body length
	if len(input.Body) == 0 {
		return nil, fmt.Errorf("body must not be empty")
	}
	if len(input.Body) > 10000 {
		return nil, fmt.Errorf("body must not exceed 10000 characters")
	}

	// Verify the showcase repo exists before allowing thread creation
	_, err := s.showcaseRepo.GetByID(ctx, input.RepoID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	now := time.Now()
	thread := &domain.Thread{
		ID:             uuid.New(),
		ShowcaseRepoID: input.RepoID,
		AuthorID:       input.AuthorID,
		Title:          input.Title,
		Body:           input.Body,
		CommentCount:   0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.threadRepo.Create(ctx, thread); err != nil {
		return nil, err
	}

	// Find the repo owner and create a notification
	repo, err := s.showcaseRepo.GetByID(ctx, input.RepoID)
	if err == nil && repo.UserID != input.AuthorID {
		// Get author alias for notification message
		author, err := s.userRepo.GetByID(ctx, input.AuthorID)
		if err == nil {
			notif := &domain.Notification{
				ID:          uuid.New(),
				UserID:      repo.UserID,
				Type:        "new_thread",
				ReferenceID: thread.ID,
				Message:     fmt.Sprintf("%s created a new thread: %s", author.Alias, input.Title),
				IsRead:      false,
				CreatedAt:   now,
			}
			_ = s.notifRepo.Create(ctx, notif)
		}
	}

	return thread, nil
}

// GetThread returns a thread with its comments.
func (s *forumService) GetThread(ctx context.Context, threadID uuid.UUID) (*domain.Thread, []domain.Comment, error) {
	thread, err := s.threadRepo.GetByID(ctx, threadID)
	if err != nil {
		return nil, nil, err
	}

	comments, err := s.commentRepo.GetByThreadID(ctx, threadID, time.Time{}, 1000)
	if err != nil {
		return nil, nil, err
	}

	return thread, comments, nil
}

// CreateComment validates input, stores the comment, increments count (eventual consistency), and notifies.
func (s *forumService) CreateComment(ctx context.Context, input domain.CreateCommentInput) (*domain.Comment, error) {
	// Sanitize: strip HTML tags from user input to prevent stored XSS
	input.Body = sanitize.StripHTMLPreserveWhitespace(input.Body)

	// Validate body length
	if len(input.Body) == 0 {
		return nil, fmt.Errorf("body must not be empty")
	}
	if len(input.Body) > 10000 {
		return nil, fmt.Errorf("body must not exceed 10000 characters")
	}

	now := time.Now()
	comment := &domain.Comment{
		ID:        uuid.New(),
		ThreadID:  input.ThreadID,
		AuthorID:  input.AuthorID,
		ParentID:  input.ParentID,
		Body:      input.Body,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Attempt to increment comment count - eventual consistency
	if err := s.threadRepo.IncrementCommentCount(ctx, input.ThreadID); err != nil {
		s.logger.Warn("failed to increment comment count",
			zap.String("thread_id", input.ThreadID.String()),
			zap.Error(err),
		)
	}

	// Create notification for thread author
	thread, err := s.threadRepo.GetByID(ctx, input.ThreadID)
	if err == nil && thread.AuthorID != input.AuthorID {
		author, err := s.userRepo.GetByID(ctx, input.AuthorID)
		if err == nil {
			notif := &domain.Notification{
				ID:          uuid.New(),
				UserID:      thread.AuthorID,
				Type:        "new_reply",
				ReferenceID: comment.ID,
				Message:     fmt.Sprintf("%s replied to your thread: %s", author.Alias, thread.Title),
				IsRead:      false,
				CreatedAt:   now,
			}
			_ = s.notifRepo.Create(ctx, notif)
		}
	}

	return comment, nil
}

// ListNotifications returns notifications for a user.
func (s *forumService) ListNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error) {
	return s.notifRepo.ListByUserID(ctx, userID, 50)
}

// MarkNotificationRead marks a notification as read for the given user.
func (s *forumService) MarkNotificationRead(ctx context.Context, userID uuid.UUID, notifID uuid.UUID) error {
	return s.notifRepo.MarkRead(ctx, userID, notifID)
}
