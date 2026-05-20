package handler

import (
	"net/http"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/upload"
	"github.com/google/uuid"
)

// BannerHandler handles profile banner upload/delete operations.
type BannerHandler struct {
	userRepo domain.UserRepository
	storage  *upload.Storage
}

// NewBannerHandler creates a new banner handler.
func NewBannerHandler(userRepo domain.UserRepository, storage *upload.Storage) *BannerHandler {
	return &BannerHandler{
		userRepo: userRepo,
		storage:  storage,
	}
}

// HandleUploadBanner handles POST /api/profile/banner (auth required, multipart/form-data).
//
// Security measures:
// - Auth required (JWT middleware applied upstream)
// - Request body limited to 10MB + overhead via http.MaxBytesReader
// - MIME type validated by inspecting actual file bytes (not trusting headers/extension)
// - Cryptographically random filename prevents path traversal and enumeration
// - Old banner file is immediately deleted from disk on replacement
// - File stored with restricted OS permissions (0640)
// - Multipart temp files cleaned up after processing
func (h *BannerHandler) HandleUploadBanner(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Security: Limit total request body to 10MB + small overhead for multipart headers
	r.Body = http.MaxBytesReader(w, r.Body, upload.MaxBannerSize+4096)

	// Parse multipart form — buffer up to 2MB in memory, rest goes to temp file
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		RespondError(w, http.StatusBadRequest, "file too large or invalid form data (max 10MB)")
		return
	}
	defer r.MultipartForm.RemoveAll() // Security: clean up temp files

	file, header, err := r.FormFile("banner")
	if err != nil {
		RespondError(w, http.StatusBadRequest, "banner file is required (form field: 'banner')")
		return
	}
	defer file.Close()

	// Security: Reject empty uploads
	if header.Size == 0 {
		RespondError(w, http.StatusBadRequest, "uploaded file is empty")
		return
	}

	// Security: Reject oversized files (defense in depth, MaxBytesReader also enforces)
	if header.Size > upload.MaxBannerSize {
		RespondError(w, http.StatusRequestEntityTooLarge, "file exceeds maximum size of 10MB")
		return
	}

	// Validate file content and store securely
	urlPath, err := h.storage.ValidateAndStore(file, header.Header.Get("Content-Type"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Replace the user's banner (delete old file, set new URL in DB)
	if err := h.replaceBanner(r, claims.UserID, urlPath); err != nil {
		// Clean up the just-stored file if DB update fails
		_ = h.storage.Delete(urlPath)
		RespondError(w, http.StatusInternalServerError, "failed to update profile banner")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"banner_url": urlPath,
		"message":    "banner uploaded successfully",
	})
}

// HandleDeleteBanner handles DELETE /api/profile/banner (auth required).
// Removes the current banner file and resets the profile to the default gradient.
func (h *BannerHandler) HandleDeleteBanner(w http.ResponseWriter, r *http.Request) {
	claims, ok := domain.GetUserClaims(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), claims.UserID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	// Delete old banner file from disk immediately
	if user.BannerURL != "" {
		_ = h.storage.Delete(user.BannerURL)
	}

	// Clear banner URL in database
	user.BannerURL = ""
	user.UpdatedAt = time.Now()
	if err := h.userRepo.Update(r.Context(), user); err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{
		"message": "banner removed",
	})
}

// replaceBanner fetches the user, deletes the old banner file, and updates the DB.
func (h *BannerHandler) replaceBanner(r *http.Request, userID uuid.UUID, newURL string) error {
	ctx := r.Context()

	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Delete old banner file from storage immediately
	if user.BannerURL != "" {
		_ = h.storage.Delete(user.BannerURL)
	}

	// Update user record with new banner URL
	user.BannerURL = newURL
	user.UpdatedAt = time.Now()
	return h.userRepo.Update(ctx, user)
}
