package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// mockThreadRepo is a mock implementation of domain.ThreadRepository.
type mockThreadRepo struct {
	threads              []domain.Thread
	incrementCountErr    error
	incrementCountCalled bool
}

func newMockThreadRepo() *mockThreadRepo {
	return &mockThreadRepo{threads: []domain.Thread{}}
}

func (m *mockThreadRepo) Create(_ context.Context, thread *domain.Thread) error {
	m.threads = append(m.threads, *thread)
	return nil
}

func (m *mockThreadRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Thread, error) {
	for i := range m.threads {
		if m.threads[i].ID == id {
			return &m.threads[i], nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockThreadRepo) ListByRepoID(_ context.Context, repoID uuid.UUID, cursor time.Time, limit int) ([]domain.Thread, error) {
	var result []domain.Thread
	for _, t := range m.threads {
		if t.ShowcaseRepoID == repoID && t.CreatedAt.Before(cursor) {
			result = append(result, t)
		}
	}
	if len(result) > limit {
		result = result[:limit]
	}
	if result == nil {
		result = []domain.Thread{}
	}
	return result, nil
}

func (m *mockThreadRepo) IncrementCommentCount(_ context.Context, _ uuid.UUID) error {
	m.incrementCountCalled = true
	return m.incrementCountErr
}

// mockCommentRepo is a mock implementation of domain.CommentRepository.
type mockCommentRepo struct {
	comments []domain.Comment
}

func newMockCommentRepo() *mockCommentRepo {
	return &mockCommentRepo{comments: []domain.Comment{}}
}

func (m *mockCommentRepo) Create(_ context.Context, comment *domain.Comment) error {
	m.comments = append(m.comments, *comment)
	return nil
}

func (m *mockCommentRepo) GetByThreadID(_ context.Context, threadID uuid.UUID, _ time.Time, limit int) ([]domain.Comment, error) {
	var result []domain.Comment
	for _, c := range m.comments {
		if c.ThreadID == threadID {
			result = append(result, c)
		}
	}
	if len(result) > limit {
		result = result[:limit]
	}
	if result == nil {
		result = []domain.Comment{}
	}
	return result, nil
}

// mockNotifRepo is a mock implementation of domain.NotificationRepository.
type mockNotifRepo struct {
	notifications []domain.Notification
}

func newMockNotifRepo() *mockNotifRepo {
	return &mockNotifRepo{notifications: []domain.Notification{}}
}

func (m *mockNotifRepo) Create(_ context.Context, notif *domain.Notification) error {
	m.notifications = append(m.notifications, *notif)
	return nil
}

func (m *mockNotifRepo) ListByUserID(_ context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	var result []domain.Notification
	for _, n := range m.notifications {
		if n.UserID == userID {
			result = append(result, n)
		}
	}
	if len(result) > limit {
		result = result[:limit]
	}
	if result == nil {
		result = []domain.Notification{}
	}
	return result, nil
}

func (m *mockNotifRepo) MarkRead(_ context.Context, userID uuid.UUID, notifID uuid.UUID) error {
	for i := range m.notifications {
		if m.notifications[i].ID == notifID && m.notifications[i].UserID == userID {
			m.notifications[i].IsRead = true
			return nil
		}
	}
	return nil
}

func setupForumService() (*forumService, *mockThreadRepo, *mockCommentRepo, *mockNotifRepo, *mockShowcaseRepo, *mockUserRepo) {
	threadRepo := newMockThreadRepo()
	commentRepo := newMockCommentRepo()
	notifRepo := newMockNotifRepo()
	showcaseRepo := newMockShowcaseRepo()
	userRepo := newMockUserRepo()
	logger := zap.NewNop()

	svc := NewForumService(threadRepo, commentRepo, notifRepo, showcaseRepo, userRepo, logger).(*forumService)
	return svc, threadRepo, commentRepo, notifRepo, showcaseRepo, userRepo
}

func TestCreateThread_ValidatesTitleTooShort(t *testing.T) {
	svc, _, _, _, _, _ := setupForumService()

	input := domain.CreateThreadInput{
		RepoID:   uuid.New(),
		AuthorID: uuid.New(),
		Title:    "Hi",
		Body:     "This is a valid body text.",
	}

	_, err := svc.CreateThread(context.Background(), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least 5 characters")
}

func TestCreateThread_ValidatesTitleTooLong(t *testing.T) {
	svc, _, _, _, _, _ := setupForumService()

	input := domain.CreateThreadInput{
		RepoID:   uuid.New(),
		AuthorID: uuid.New(),
		Title:    strings.Repeat("a", 256),
		Body:     "This is a valid body text.",
	}

	_, err := svc.CreateThread(context.Background(), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not exceed 255 characters")
}

func TestCreateThread_ValidatesBodyEmpty(t *testing.T) {
	svc, _, _, _, _, _ := setupForumService()

	input := domain.CreateThreadInput{
		RepoID:   uuid.New(),
		AuthorID: uuid.New(),
		Title:    "Valid Title",
		Body:     "",
	}

	_, err := svc.CreateThread(context.Background(), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "body must not be empty")
}

func TestCreateThread_ValidatesBodyTooLong(t *testing.T) {
	svc, _, _, _, _, _ := setupForumService()

	input := domain.CreateThreadInput{
		RepoID:   uuid.New(),
		AuthorID: uuid.New(),
		Title:    "Valid Title",
		Body:     strings.Repeat("x", 10001),
	}

	_, err := svc.CreateThread(context.Background(), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not exceed 10000 characters")
}

func TestCreateThread_ReturnsNotFoundWhenRepoDoesNotExist(t *testing.T) {
	svc, _, _, _, _, _ := setupForumService()

	// Use a RepoID that does not match any showcase repo in the mock
	nonExistentRepoID := uuid.New()

	input := domain.CreateThreadInput{
		RepoID:   nonExistentRepoID,
		AuthorID: uuid.New(),
		Title:    "Valid Thread Title",
		Body:     "This is a valid body text for the thread.",
	}

	thread, err := svc.CreateThread(context.Background(), input)
	require.Error(t, err)
	assert.Nil(t, thread)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestCreateThread_CreatesNotificationForRepoOwner(t *testing.T) {
	svc, _, _, notifRepo, showcaseRepo, userRepo := setupForumService()

	ownerID := uuid.New()
	authorID := uuid.New()
	repoID := uuid.New()

	// Set up repo owner
	userRepo.users["owner"] = &domain.User{
		ID:    ownerID,
		NIM:   "owner",
		Alias: "repo_owner",
	}
	// Set up author
	userRepo.users["author"] = &domain.User{
		ID:    authorID,
		NIM:   "author",
		Alias: "thread_author",
	}
	// Set up showcase repo
	showcaseRepo.repos = append(showcaseRepo.repos, domain.ShowcaseRepo{
		ID:     repoID,
		UserID: ownerID,
	})

	input := domain.CreateThreadInput{
		RepoID:   repoID,
		AuthorID: authorID,
		Title:    "Hello World Thread",
		Body:     "This is a valid body.",
	}

	thread, err := svc.CreateThread(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, thread)

	// Verify notification was created for repo owner
	require.Len(t, notifRepo.notifications, 1)
	notif := notifRepo.notifications[0]
	assert.Equal(t, ownerID, notif.UserID)
	assert.Equal(t, "new_thread", notif.Type)
	assert.Equal(t, thread.ID, notif.ReferenceID)
	assert.Contains(t, notif.Message, "thread_author")
	assert.Contains(t, notif.Message, "Hello World Thread")
}

func TestCreateThread_NoNotificationWhenAuthorIsOwner(t *testing.T) {
	svc, _, _, notifRepo, showcaseRepo, userRepo := setupForumService()

	ownerID := uuid.New()
	repoID := uuid.New()

	// Same user is both owner and author
	userRepo.users["owner"] = &domain.User{
		ID:    ownerID,
		NIM:   "owner",
		Alias: "the_owner",
	}
	showcaseRepo.repos = append(showcaseRepo.repos, domain.ShowcaseRepo{
		ID:     repoID,
		UserID: ownerID,
	})

	input := domain.CreateThreadInput{
		RepoID:   repoID,
		AuthorID: ownerID,
		Title:    "My Own Thread",
		Body:     "This is posted by the owner themselves.",
	}

	_, err := svc.CreateThread(context.Background(), input)
	require.NoError(t, err)

	// No notification should be created when author is the repo owner
	assert.Len(t, notifRepo.notifications, 0)
}

func TestCreateComment_ValidatesBodyEmpty(t *testing.T) {
	svc, _, _, _, _, _ := setupForumService()

	input := domain.CreateCommentInput{
		ThreadID: uuid.New(),
		AuthorID: uuid.New(),
		Body:     "",
	}

	_, err := svc.CreateComment(context.Background(), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "body must not be empty")
}

func TestCreateComment_ValidatesBodyTooLong(t *testing.T) {
	svc, _, _, _, _, _ := setupForumService()

	input := domain.CreateCommentInput{
		ThreadID: uuid.New(),
		AuthorID: uuid.New(),
		Body:     strings.Repeat("y", 10001),
	}

	_, err := svc.CreateComment(context.Background(), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not exceed 10000 characters")
}

func TestCreateComment_CreatesNotificationForThreadAuthor(t *testing.T) {
	svc, threadRepo, _, notifRepo, _, userRepo := setupForumService()

	threadAuthorID := uuid.New()
	commentAuthorID := uuid.New()
	threadID := uuid.New()

	// Set up users
	userRepo.users["thread_author"] = &domain.User{
		ID:    threadAuthorID,
		NIM:   "thread_author",
		Alias: "thread_author_alias",
	}
	userRepo.users["comment_author"] = &domain.User{
		ID:    commentAuthorID,
		NIM:   "comment_author",
		Alias: "commenter",
	}

	// Create a thread
	threadRepo.threads = append(threadRepo.threads, domain.Thread{
		ID:       threadID,
		AuthorID: threadAuthorID,
		Title:    "Test Thread",
	})

	input := domain.CreateCommentInput{
		ThreadID: threadID,
		AuthorID: commentAuthorID,
		Body:     "This is a comment.",
	}

	comment, err := svc.CreateComment(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, comment)

	// Verify notification was created for thread author
	require.Len(t, notifRepo.notifications, 1)
	notif := notifRepo.notifications[0]
	assert.Equal(t, threadAuthorID, notif.UserID)
	assert.Equal(t, "new_reply", notif.Type)
	assert.Equal(t, comment.ID, notif.ReferenceID)
	assert.Contains(t, notif.Message, "commenter")
	assert.Contains(t, notif.Message, "Test Thread")
}

func TestCreateComment_EventualConsistency_CommentStoredEvenIfCountFails(t *testing.T) {
	svc, threadRepo, commentRepo, _, _, userRepo := setupForumService()

	threadAuthorID := uuid.New()
	commentAuthorID := uuid.New()
	threadID := uuid.New()

	userRepo.users["ta"] = &domain.User{
		ID:    threadAuthorID,
		NIM:   "ta",
		Alias: "ta_alias",
	}
	userRepo.users["ca"] = &domain.User{
		ID:    commentAuthorID,
		NIM:   "ca",
		Alias: "ca_alias",
	}

	threadRepo.threads = append(threadRepo.threads, domain.Thread{
		ID:       threadID,
		AuthorID: threadAuthorID,
		Title:    "Thread for Count Test",
	})

	// Make IncrementCommentCount fail
	threadRepo.incrementCountErr = errors.New("database connection lost")

	input := domain.CreateCommentInput{
		ThreadID: threadID,
		AuthorID: commentAuthorID,
		Body:     "Comment that should still be stored.",
	}

	comment, err := svc.CreateComment(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, comment)

	// Comment should be stored even though count increment failed
	assert.Len(t, commentRepo.comments, 1)
	assert.Equal(t, "Comment that should still be stored.", commentRepo.comments[0].Body)

	// Verify IncrementCommentCount was called
	assert.True(t, threadRepo.incrementCountCalled)
}

func TestCreateComment_NoNotificationWhenCommentAuthorIsThreadAuthor(t *testing.T) {
	svc, threadRepo, _, notifRepo, _, userRepo := setupForumService()

	authorID := uuid.New()
	threadID := uuid.New()

	userRepo.users["author"] = &domain.User{
		ID:    authorID,
		NIM:   "author",
		Alias: "self_replier",
	}

	threadRepo.threads = append(threadRepo.threads, domain.Thread{
		ID:       threadID,
		AuthorID: authorID,
		Title:    "My Thread",
	})

	input := domain.CreateCommentInput{
		ThreadID: threadID,
		AuthorID: authorID,
		Body:     "Replying to my own thread.",
	}

	_, err := svc.CreateComment(context.Background(), input)
	require.NoError(t, err)

	// No notification when commenting on own thread
	assert.Len(t, notifRepo.notifications, 0)
}

func TestNotificationTargeting_OnlyRepoOwnerGetsThreadNotification(t *testing.T) {
	svc, _, _, notifRepo, showcaseRepo, userRepo := setupForumService()

	ownerID := uuid.New()
	authorID := uuid.New()
	bystander := uuid.New()
	repoID := uuid.New()

	userRepo.users["owner"] = &domain.User{ID: ownerID, NIM: "owner", Alias: "owner_alias"}
	userRepo.users["author"] = &domain.User{ID: authorID, NIM: "author", Alias: "author_alias"}
	userRepo.users["bystander"] = &domain.User{ID: bystander, NIM: "bystander", Alias: "bystander_alias"}

	showcaseRepo.repos = append(showcaseRepo.repos, domain.ShowcaseRepo{
		ID:     repoID,
		UserID: ownerID,
	})

	input := domain.CreateThreadInput{
		RepoID:   repoID,
		AuthorID: authorID,
		Title:    "Thread for Targeting Test",
		Body:     "Testing notification targeting.",
	}

	_, err := svc.CreateThread(context.Background(), input)
	require.NoError(t, err)

	// Only owner should get the notification, not the bystander
	require.Len(t, notifRepo.notifications, 1)
	assert.Equal(t, ownerID, notifRepo.notifications[0].UserID)
}

func TestNotificationTargeting_OnlyThreadAuthorGetsCommentNotification(t *testing.T) {
	svc, threadRepo, _, notifRepo, _, userRepo := setupForumService()

	threadAuthorID := uuid.New()
	commenterID := uuid.New()
	bystander := uuid.New()
	threadID := uuid.New()

	userRepo.users["ta"] = &domain.User{ID: threadAuthorID, NIM: "ta", Alias: "thread_auth"}
	userRepo.users["cm"] = &domain.User{ID: commenterID, NIM: "cm", Alias: "commenter_alias"}
	userRepo.users["by"] = &domain.User{ID: bystander, NIM: "by", Alias: "bystander"}

	threadRepo.threads = append(threadRepo.threads, domain.Thread{
		ID:       threadID,
		AuthorID: threadAuthorID,
		Title:    "Thread for Comment Targeting",
	})

	input := domain.CreateCommentInput{
		ThreadID: threadID,
		AuthorID: commenterID,
		Body:     "This is a comment.",
	}

	_, err := svc.CreateComment(context.Background(), input)
	require.NoError(t, err)

	// Only thread author should get the notification
	require.Len(t, notifRepo.notifications, 1)
	assert.Equal(t, threadAuthorID, notifRepo.notifications[0].UserID)
}

func TestListNotifications_ReturnsUserNotifications(t *testing.T) {
	svc, _, _, notifRepo, _, _ := setupForumService()

	userID := uuid.New()
	otherUserID := uuid.New()

	notifRepo.notifications = []domain.Notification{
		{ID: uuid.New(), UserID: userID, Type: "new_thread", Message: "msg1", CreatedAt: time.Now()},
		{ID: uuid.New(), UserID: otherUserID, Type: "new_reply", Message: "msg2", CreatedAt: time.Now()},
		{ID: uuid.New(), UserID: userID, Type: "new_reply", Message: "msg3", CreatedAt: time.Now()},
	}

	result, err := svc.ListNotifications(context.Background(), userID)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestMarkNotificationRead_UpdatesCorrectNotification(t *testing.T) {
	svc, _, _, notifRepo, _, _ := setupForumService()

	userID := uuid.New()
	notifID := uuid.New()

	notifRepo.notifications = []domain.Notification{
		{ID: notifID, UserID: userID, Type: "new_thread", IsRead: false, CreatedAt: time.Now()},
	}

	err := svc.MarkNotificationRead(context.Background(), userID, notifID)
	require.NoError(t, err)
	assert.True(t, notifRepo.notifications[0].IsRead)
}
