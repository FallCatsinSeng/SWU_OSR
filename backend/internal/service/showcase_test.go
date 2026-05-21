package service

import (
	"context"
	"testing"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetShowcase_ValidatesMaxLimit(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	ghSvc := &mockGitHubService{}

	userID := uuid.New()
	key := []byte("01234567890123456789012345678901") // 32 bytes
	encToken, _ := Encrypt("gh-token", key)

	user := &domain.User{
		ID:          userID,
		NIM:         "2021001",
		Alias:       "testuser",
		GitHubToken: encToken,
	}
	userRepo.users["2021001"] = user

	svc := NewShowcaseService(showcaseRepo, userRepo, ghSvc, key)

	// Create 21 selections - should fail
	selections := make([]domain.ShowcaseSelection, 21)
	for i := range selections {
		selections[i] = domain.ShowcaseSelection{
			RepoID:   int64(i + 1),
			RepoName: "repo",
			FullName: "user/repo",
			Tag:      domain.TagCoursework,
		}
	}

	err := svc.SetShowcase(context.Background(), userID, selections)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maximum 20")
}

func TestSetShowcase_ValidatesAcademicTags(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	ghSvc := &mockGitHubService{}

	userID := uuid.New()
	key := []byte("01234567890123456789012345678901")
	encToken, _ := Encrypt("gh-token", key)

	user := &domain.User{
		ID:          userID,
		NIM:         "2021001",
		Alias:       "testuser",
		GitHubToken: encToken,
	}
	userRepo.users["2021001"] = user

	svc := NewShowcaseService(showcaseRepo, userRepo, ghSvc, key)

	// Invalid tag
	selections := []domain.ShowcaseSelection{
		{
			RepoID:   1,
			RepoName: "repo1",
			FullName: "user/repo1",
			Tag:      domain.AcademicTag("invalid_tag"),
		},
	}

	err := svc.SetShowcase(context.Background(), userID, selections)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid academic tag")
}

func TestSetShowcase_ValidTags(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	ghSvc := &mockGitHubService{}

	userID := uuid.New()
	key := []byte("01234567890123456789012345678901")
	encToken, _ := Encrypt("gh-token", key)

	user := &domain.User{
		ID:          userID,
		NIM:         "2021001",
		Alias:       "testuser",
		GitHubToken: encToken,
	}
	userRepo.users["2021001"] = user

	svc := NewShowcaseService(showcaseRepo, userRepo, ghSvc, key)

	validTags := []domain.AcademicTag{
		domain.TagCoursework,
		domain.TagThesis,
		domain.TagHackathon,
		domain.TagPersonalResearch,
		domain.TagTeamProject,
	}

	for _, tag := range validTags {
		// Reset showcase repo to avoid accumulation
		showcaseRepo.repos = nil

		selections := []domain.ShowcaseSelection{
			{
				RepoID:   1,
				RepoName: "repo1",
				FullName: "user/repo1",
				Tag:      tag,
			},
		}

		err := svc.SetShowcase(context.Background(), userID, selections)
		assert.NoError(t, err, "tag %s should be valid", tag)
	}
}

func TestSetShowcase_Max20Succeeds(t *testing.T) {
	userRepo := newMockUserRepo()
	showcaseRepo := newMockShowcaseRepo()
	ghSvc := &mockGitHubService{}

	userID := uuid.New()
	key := []byte("01234567890123456789012345678901")
	encToken, _ := Encrypt("gh-token", key)

	user := &domain.User{
		ID:          userID,
		NIM:         "2021001",
		Alias:       "testuser",
		GitHubToken: encToken,
	}
	userRepo.users["2021001"] = user

	svc := NewShowcaseService(showcaseRepo, userRepo, ghSvc, key)

	// Exactly 20 selections - should succeed
	selections := make([]domain.ShowcaseSelection, 20)
	for i := range selections {
		selections[i] = domain.ShowcaseSelection{
			RepoID:   int64(i + 1),
			RepoName: "repo",
			FullName: "user/repo",
			Tag:      domain.TagCoursework,
		}
	}

	err := svc.SetShowcase(context.Background(), userID, selections)
	assert.NoError(t, err)
}
