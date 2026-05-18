package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testWebhookSecret = "test-webhook-secret"

func computeSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestProcessWebhook_ValidSignature(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	activityRepo := newMockActivityRepo()

	userID := uuid.New()
	user := &domain.User{
		ID:             userID,
		NIM:            "2021001",
		Alias:          "testuser",
		GitHubUsername: "ghuser",
	}
	userRepo.users["2021001"] = user

	repoID := uuid.New()
	showcaseRepo.repos = append(showcaseRepo.repos, domain.ShowcaseRepo{
		ID:           repoID,
		UserID:       userID,
		RepoFullName: "ghuser/myrepo",
	})

	svc := NewAggregatorService(activityRepo, userRepo, showcaseRepo, testWebhookSecret)

	payload := []byte(`{"repository":{"full_name":"ghuser/myrepo"},"pusher":{"name":"ghuser"},"ref":"refs/heads/main","commits":[{"id":"abc123","message":"test","author":{"name":"ghuser"}}]}`)
	signature := computeSignature(payload, testWebhookSecret)

	err := svc.ProcessWebhook(context.Background(), payload, signature, "push", "delivery-1")
	require.NoError(t, err)
	assert.Len(t, activityRepo.logs, 1)
}

func TestProcessWebhook_InvalidSignature(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	activityRepo := newMockActivityRepo()

	svc := NewAggregatorService(activityRepo, userRepo, showcaseRepo, testWebhookSecret)

	payload := []byte(`{"test": true}`)
	err := svc.ProcessWebhook(context.Background(), payload, "sha256=invalid", "push", "delivery-1")
	assert.ErrorIs(t, err, domain.ErrInvalidSignature)
}

func TestProcessWebhook_PushEvent(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	activityRepo := newMockActivityRepo()

	userID := uuid.New()
	userRepo.users["2021001"] = &domain.User{
		ID: userID, NIM: "2021001", Alias: "pushuser", GitHubUsername: "pushuser",
	}
	showcaseRepo.repos = append(showcaseRepo.repos, domain.ShowcaseRepo{
		ID: uuid.New(), UserID: userID, RepoFullName: "pushuser/repo",
	})

	svc := NewAggregatorService(activityRepo, userRepo, showcaseRepo, testWebhookSecret)

	payload := []byte(`{"repository":{"full_name":"pushuser/repo"},"pusher":{"name":"pushuser"},"ref":"refs/heads/main","commits":[{"id":"1","message":"first","author":{"name":"pushuser"}},{"id":"2","message":"second","author":{"name":"pushuser"}}]}`)
	signature := computeSignature(payload, testWebhookSecret)

	err := svc.ProcessWebhook(context.Background(), payload, signature, "push", "delivery-push-1")
	require.NoError(t, err)
	require.Len(t, activityRepo.logs, 1)
	assert.Equal(t, domain.EventPush, activityRepo.logs[0].EventType)
	assert.Contains(t, activityRepo.logs[0].Summary, "2 commit")
}

func TestProcessWebhook_PREvent(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	activityRepo := newMockActivityRepo()

	userID := uuid.New()
	userRepo.users["2021001"] = &domain.User{
		ID: userID, NIM: "2021001", Alias: "pruser", GitHubUsername: "pruser",
	}
	showcaseRepo.repos = append(showcaseRepo.repos, domain.ShowcaseRepo{
		ID: uuid.New(), UserID: userID, RepoFullName: "pruser/repo",
	})

	svc := NewAggregatorService(activityRepo, userRepo, showcaseRepo, testWebhookSecret)

	payload := []byte(`{"action":"opened","number":42,"repository":{"full_name":"pruser/repo"},"pull_request":{"title":"Add feature X","user":{"login":"pruser"}}}`)
	signature := computeSignature(payload, testWebhookSecret)

	err := svc.ProcessWebhook(context.Background(), payload, signature, "pull_request", "delivery-pr-1")
	require.NoError(t, err)
	require.Len(t, activityRepo.logs, 1)
	assert.Equal(t, domain.EventPR, activityRepo.logs[0].EventType)
	assert.Contains(t, activityRepo.logs[0].Summary, "#42")
}

func TestProcessWebhook_ReleaseEvent(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	activityRepo := newMockActivityRepo()

	userID := uuid.New()
	userRepo.users["2021001"] = &domain.User{
		ID: userID, NIM: "2021001", Alias: "reluser", GitHubUsername: "reluser",
	}
	showcaseRepo.repos = append(showcaseRepo.repos, domain.ShowcaseRepo{
		ID: uuid.New(), UserID: userID, RepoFullName: "reluser/repo",
	})

	svc := NewAggregatorService(activityRepo, userRepo, showcaseRepo, testWebhookSecret)

	payload := []byte(`{"repository":{"full_name":"reluser/repo"},"release":{"tag_name":"v1.0.0","name":"Version 1.0","author":{"login":"reluser"}}}`)
	signature := computeSignature(payload, testWebhookSecret)

	err := svc.ProcessWebhook(context.Background(), payload, signature, "release", "delivery-rel-1")
	require.NoError(t, err)
	require.Len(t, activityRepo.logs, 1)
	assert.Equal(t, domain.EventRelease, activityRepo.logs[0].EventType)
	assert.Contains(t, activityRepo.logs[0].Summary, "v1.0.0")
}

func TestProcessWebhook_Idempotency(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	activityRepo := newMockActivityRepo()

	userID := uuid.New()
	userRepo.users["2021001"] = &domain.User{
		ID: userID, NIM: "2021001", Alias: "idempuser", GitHubUsername: "idempuser",
	}
	showcaseRepo.repos = append(showcaseRepo.repos, domain.ShowcaseRepo{
		ID: uuid.New(), UserID: userID, RepoFullName: "idempuser/repo",
	})

	svc := NewAggregatorService(activityRepo, userRepo, showcaseRepo, testWebhookSecret)

	payload := []byte(`{"repository":{"full_name":"idempuser/repo"},"pusher":{"name":"idempuser"},"ref":"refs/heads/main","commits":[{"id":"abc","message":"test","author":{"name":"idempuser"}}]}`)
	signature := computeSignature(payload, testWebhookSecret)

	err := svc.ProcessWebhook(context.Background(), payload, signature, "push", "same-delivery-id")
	require.NoError(t, err)
	assert.Len(t, activityRepo.logs, 1)

	// Process same delivery again
	err = svc.ProcessWebhook(context.Background(), payload, signature, "push", "same-delivery-id")
	require.NoError(t, err)
	assert.Len(t, activityRepo.logs, 1) // Still 1
}

func TestGetActivityFeed_LimitClamping(t *testing.T) {
	activityRepo := newMockActivityRepo()
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()

	svc := NewAggregatorService(activityRepo, userRepo, showcaseRepo, testWebhookSecret)

	result, err := svc.GetActivityFeed(context.Background(), domain.FeedParams{Limit: 0})
	require.NoError(t, err)
	assert.NotNil(t, result)

	result, err = svc.GetActivityFeed(context.Background(), domain.FeedParams{Limit: -5})
	require.NoError(t, err)
	assert.NotNil(t, result)

	result, err = svc.GetActivityFeed(context.Background(), domain.FeedParams{Limit: 100})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCursorEncodeDecode(t *testing.T) {
	now := time.Now()
	encoded := encodeCursor(now)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)
	parsedTime, err := time.Parse(time.RFC3339Nano, string(decoded))
	require.NoError(t, err)
	assert.WithinDuration(t, now, parsedTime, time.Microsecond)
}

func TestClampLimit(t *testing.T) {
	assert.Equal(t, 20, clampLimit(0))
	assert.Equal(t, 20, clampLimit(-1))
	assert.Equal(t, 1, clampLimit(1))
	assert.Equal(t, 50, clampLimit(50))
	assert.Equal(t, 50, clampLimit(51))
	assert.Equal(t, 25, clampLimit(25))
}

func TestDecodeCursor_Empty(t *testing.T) {
	result := decodeCursor("")
	assert.True(t, result.After(time.Now().Add(-time.Second)))
}

func TestDecodeCursor_Invalid(t *testing.T) {
	result := decodeCursor("not-valid-base64")
	assert.True(t, result.After(time.Now().Add(-time.Second)))
}

func TestVerifySignature(t *testing.T) {
	svc := &aggregatorService{webhookSecret: []byte(testWebhookSecret)}
	payload := []byte(`{"test": "data"}`)
	validSig := computeSignature(payload, testWebhookSecret)
	assert.True(t, svc.verifySignature(payload, validSig))
	assert.False(t, svc.verifySignature(payload, "sha256=wrong"))
	assert.False(t, svc.verifySignature(payload, ""))
}

func TestParsePushEvent(t *testing.T) {
	svc := &aggregatorService{}
	payload := []byte(`{"repository":{"full_name":"user/repo"},"pusher":{"name":"user"},"ref":"refs/heads/feature-branch","commits":[{"id":"1","message":"commit 1","author":{"name":"user"}},{"id":"2","message":"commit 2","author":{"name":"user"}},{"id":"3","message":"commit 3","author":{"name":"user"}}]}`)
	repoFullName, username, summary, metadata := svc.parsePushEvent(payload)
	assert.Equal(t, "user/repo", repoFullName)
	assert.Equal(t, "user", username)
	assert.Contains(t, summary, "3 commit")
	assert.Contains(t, summary, "feature-branch")
	var meta map[string]interface{}
	err := json.Unmarshal(metadata, &meta)
	require.NoError(t, err)
	assert.Equal(t, float64(3), meta["commit_count"])
}

func TestParsePREvent(t *testing.T) {
	svc := &aggregatorService{}
	payload := []byte(`{"action":"closed","number":10,"repository":{"full_name":"user/repo"},"pull_request":{"title":"Fix bug","user":{"login":"developer"}}}`)
	repoFullName, username, summary, _ := svc.parsePREvent(payload)
	assert.Equal(t, "user/repo", repoFullName)
	assert.Equal(t, "developer", username)
	assert.Contains(t, summary, "#10")
	assert.Contains(t, summary, "closed")
}

func TestParseReleaseEvent(t *testing.T) {
	svc := &aggregatorService{}
	payload := []byte(`{"repository":{"full_name":"org/project"},"release":{"tag_name":"v2.0.0","name":"Major Release","author":{"login":"maintainer"}}}`)
	repoFullName, username, summary, _ := svc.parseReleaseEvent(payload)
	assert.Equal(t, "org/project", repoFullName)
	assert.Equal(t, "maintainer", username)
	assert.Contains(t, summary, "v2.0.0")
	assert.Contains(t, summary, "Major Release")
}
