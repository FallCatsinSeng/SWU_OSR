package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// OAuthToken holds the token received from GitHub OAuth exchange.
type OAuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GitHubUser represents a GitHub user profile.
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

// Repository represents a GitHub repository.
type Repository struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Language    string `json:"language"`
}

// Service defines the GitHub OAuth and API interface.
type Service interface {
	GetAuthorizationURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*OAuthToken, error)
	GetUser(ctx context.Context, token string) (*GitHubUser, error)
	ListRepos(ctx context.Context, token string) ([]Repository, error)
	RegisterWebhook(ctx context.Context, token, owner, repo, webhookURL, secret string) (int64, error)
	RemoveWebhook(ctx context.Context, token, owner, repo string, hookID int64) error
}

// service is the concrete implementation.
type service struct {
	clientID     string
	clientSecret string
	redirectURI  string
	httpClient   *http.Client
}

// NewService creates a new GitHub service.
func NewService(clientID, clientSecret, redirectURI string) Service {
	return &service{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		httpClient:   &http.Client{},
	}
}

// GetAuthorizationURL constructs the GitHub OAuth authorization URL.
func (s *service) GetAuthorizationURL(state string) string {
	params := url.Values{}
	params.Set("client_id", s.clientID)
	params.Set("redirect_uri", s.redirectURI)
	params.Set("scope", "read:user repo")
	params.Set("state", state)
	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for an access token.
func (s *service) ExchangeCode(ctx context.Context, code string) (*OAuthToken, error) {
	form := url.Values{}
	form.Set("client_id", s.clientID)
	form.Set("client_secret", s.clientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", s.redirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://github.com/login/oauth/access_token",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating exchange request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchanging code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var token OAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	if token.AccessToken == "" {
		return nil, fmt.Errorf("empty access token in response")
	}

	return &token, nil
}

// GetUser retrieves the authenticated GitHub user's profile.
func (s *service) GetUser(ctx context.Context, token string) (*GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("creating user request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get user failed with status %d: %s", resp.StatusCode, string(body))
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decoding user response: %w", err)
	}

	return &user, nil
}

// ListRepos retrieves all repositories for the authenticated user (paginated).
func (s *service) ListRepos(ctx context.Context, token string) ([]Repository, error) {
	var allRepos []Repository
	page := 1

	for {
		reqURL := fmt.Sprintf("https://api.github.com/user/repos?per_page=100&page=%d", page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("creating repos request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetching repos page %d: %w", page, err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("list repos failed with status %d: %s", resp.StatusCode, string(body))
		}

		var repos []Repository
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decoding repos response: %w", err)
		}
		resp.Body.Close()

		if len(repos) == 0 {
			break
		}

		allRepos = append(allRepos, repos...)
		page++

		if len(repos) < 100 {
			break
		}
	}

	return allRepos, nil
}

// RegisterWebhook creates a webhook on the given repository.
func (s *service) RegisterWebhook(ctx context.Context, token, owner, repo, webhookURL, secret string) (int64, error) {
	payload := map[string]interface{}{
		"name":   "web",
		"active": true,
		"events": []string{"push", "pull_request"},
		"config": map[string]string{
			"url":          webhookURL,
			"content_type": "json",
			"secret":       secret,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("marshaling webhook payload: %w", err)
	}

	reqURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(string(body)))
	if err != nil {
		return 0, fmt.Errorf("creating webhook request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("registering webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("register webhook failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decoding webhook response: %w", err)
	}

	return result.ID, nil
}

// RemoveWebhook deletes a webhook from the given repository.
func (s *service) RemoveWebhook(ctx context.Context, token, owner, repo string, hookID int64) error {
	reqURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks/%d", owner, repo, hookID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating delete webhook request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("removing webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("remove webhook failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
