package siakad

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate_ValidLogin(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		// Set PHPSESSID cookie on first warmup request
		if requestCount == 1 {
			http.SetCookie(w, &http.Cookie{
				Name:  "PHPSESSID",
				Value: "test-session-abc123",
			})
		}

		// Handle login POST
		if r.URL.Path == "/login_proses.php" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Login berhasil")
			return
		}

		// Warmup pages
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	svc := NewService(server.URL, 10*time.Second)
	data, err := svc.Authenticate(context.Background(), "123456", "password")

	require.NoError(t, err)
	assert.Equal(t, "123456", data.NIM)
	assert.Equal(t, "test-session-abc123", data.SessionID)
	assert.True(t, data.IsActive)
}

func TestAuthenticate_InvalidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login_proses.php" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "tidakterdaftar")
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "PHPSESSID",
			Value: "test-session",
		})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc := NewService(server.URL, 10*time.Second)
	_, err := svc.Authenticate(context.Background(), "123456", "wrong")

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestAuthenticate_DeviceRejected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login_proses.php" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "salahdevice")
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "PHPSESSID",
			Value: "test-session",
		})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc := NewService(server.URL, 10*time.Second)
	_, err := svc.Authenticate(context.Background(), "123456", "password")

	assert.ErrorIs(t, err, domain.ErrDeviceRejected)
}

func TestAuthenticate_ServerTimeout_RetriesAndFails(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return 500 on first warmup request to trigger retries
		if strings.Contains(r.URL.Path, "swu.php") {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc := NewService(server.URL, 5*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := svc.Authenticate(ctx, "123456", "password")

	assert.ErrorIs(t, err, domain.ErrSIAKADUnavailable)
	assert.Equal(t, 3, attempts)
}
