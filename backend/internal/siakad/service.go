package siakad

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
)

// StudentData holds the information retrieved from SIAKAD after a successful login.
type StudentData struct {
	NIM       string
	FullName  string
	Major     string
	Semester  int
	IsActive  bool
	SessionID string
}

// Service defines the SIAKAD proxy interface.
type Service interface {
	Authenticate(ctx context.Context, nim, password string) (*StudentData, error)
}

// service is the concrete implementation of the SIAKAD proxy.
type service struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewService creates a new SIAKAD proxy service.
func NewService(baseURL string, timeout time.Duration) Service {
	return &service{
		baseURL: strings.TrimRight(baseURL, "/"),
		timeout: timeout,
	}
}

// warmupChain is the list of URLs to visit before authenticating.
var warmupChain = []string{
	"/swu.php",
	"/my_school.php?ada=2&sof=0&ol=0&hp=1&template=0",
	"/my_school_ok.php?benarinput=0&ada=2&sof=0&ol=0&hp=1&template=0",
	"/my_school_run.php?ada=2&sof=0&ol=0&hp=1&template=0",
	"/smart_school_biasa_2019.php",
}

// Authenticate performs SIAKAD login with warmup chain and retry logic.
func (s *service) Authenticate(ctx context.Context, nim, password string) (*StudentData, error) {
	var lastErr error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		data, err := s.doAuthenticate(ctx, nim, password)
		if err != nil {
			// Only retry on network errors or 5xx responses; do not retry credential/device errors.
			if isRetryable(err) {
				lastErr = err
				continue
			}
			return nil, err
		}
		return data, nil
	}

	return nil, fmt.Errorf("%w: %v", domain.ErrSIAKADUnavailable, lastErr)
}

func (s *service) doAuthenticate(ctx context.Context, nim, password string) (*StudentData, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, domain.ErrSessionInitFailed
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: s.timeout,
	}

	// Execute warmup chain
	for _, path := range warmupChain {
		reqURL := s.baseURL + path
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, domain.ErrSessionInitFailed
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, &retryableError{err}
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 500 {
			return nil, &retryableError{fmt.Errorf("warmup returned status %d", resp.StatusCode)}
		}
	}

	// POST login
	form := url.Values{}
	form.Set("username", nim)
	form.Set("password", password)
	form.Set("mac_addr", "02:00:00:00:00:00")

	loginURL := s.baseURL + "/login_proses.php"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, domain.ErrSessionInitFailed
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, &retryableError{err}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, &retryableError{fmt.Errorf("login returned status %d", resp.StatusCode)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &retryableError{err}
	}

	bodyStr := strings.ToLower(string(body))

	if strings.Contains(bodyStr, "tidakterdaftar") {
		return nil, domain.ErrInvalidCredentials
	}
	if strings.Contains(bodyStr, "salahdevice") {
		return nil, domain.ErrDeviceRejected
	}

	// Extract PHPSESSID from cookies
	parsedURL, _ := url.Parse(s.baseURL)
	var sessionID string
	for _, c := range jar.Cookies(parsedURL) {
		if c.Name == "PHPSESSID" {
			sessionID = c.Value
			break
		}
	}

	if sessionID == "" {
		return nil, domain.ErrSessionInitFailed
	}

	return &StudentData{
		NIM:       nim,
		FullName:  nim, // SIAKAD does not return full profile in login response; populated later
		Major:     "",
		Semester:  0,
		IsActive:  true,
		SessionID: sessionID,
	}, nil
}

// retryableError marks an error as retryable.
type retryableError struct {
	err error
}

func (e *retryableError) Error() string {
	return e.err.Error()
}

func (e *retryableError) Unwrap() error {
	return e.err
}

func isRetryable(err error) bool {
	_, ok := err.(*retryableError)
	return ok
}
