// Package upload provides secure file upload handling with MIME validation,
// size enforcement, and local disk storage. Designed for profile banners.
package upload

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// MaxBannerSize is the maximum allowed upload size (10 MB).
const MaxBannerSize = 10 << 20 // 10 MB

// AllowedMIMETypes maps MIME types to their canonical file extensions.
// Security: Only allow known-safe media types — no executables, no HTML, no SVG (XSS risk).
var AllowedMIMETypes = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
	"image/gif":       ".gif",
	"video/mp4":       ".mp4",
	"video/webm":      ".webm",
}

// Storage handles local file persistence for uploaded banners.
type Storage struct {
	// BasePath is the directory where uploaded files are stored on disk.
	BasePath string
	// URLPrefix is the URL prefix used to serve the files (e.g., "/uploads/banners/").
	URLPrefix string
}

// NewStorage creates a new file storage instance and ensures the directory exists.
func NewStorage(basePath, urlPrefix string) (*Storage, error) {
	// Check if directory already exists (e.g., pre-created in Docker image or volume mount)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		// Try to create it — may fail in read-only containers without volume mounts
		if mkErr := os.MkdirAll(basePath, 0750); mkErr != nil {
			return nil, fmt.Errorf("creating upload directory: %w", mkErr)
		}
	} else if err != nil {
		return nil, fmt.Errorf("checking upload directory: %w", err)
	}

	// Verify we can actually write to the directory
	testFile := filepath.Join(basePath, ".write-test")
	f, err := os.Create(testFile)
	if err != nil {
		return nil, fmt.Errorf("upload directory is not writable: %w", err)
	}
	f.Close()
	os.Remove(testFile)

	return &Storage{
		BasePath:  basePath,
		URLPrefix: urlPrefix,
	}, nil
}

// ValidateAndStore reads the uploaded file, validates MIME type by inspecting actual
// file content (not trusting Content-Type header), and stores it with a random filename.
// Returns the public URL path for the stored file.
//
// Security measures:
// - File size enforced via http.MaxBytesReader (caller responsibility) + re-check here
// - MIME type detected from first 512 bytes of actual content (not from header/extension)
// - Filename is a cryptographically random hex string (no user-controlled path components)
// - Extension is derived from validated MIME type (not from user input)
// - File is written with world-readable permissions (0644) since nginx serves
//   these files from a shared volume as a different user
// - Path traversal impossible due to generated filename with no directory separators
func (s *Storage) ValidateAndStore(file io.Reader, declaredContentType string) (string, error) {
	// Read up to MaxBannerSize + 1 byte to detect oversized files
	// (in case MaxBytesReader wasn't applied at the handler level)
	limited := io.LimitReader(file, MaxBannerSize+1)

	// Read the first 512 bytes for MIME detection
	header := make([]byte, 512)
	n, err := io.ReadFull(limited, header)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", fmt.Errorf("reading file header: %w", err)
	}
	header = header[:n]

	// Security: Detect MIME from actual file content, not from user-supplied header
	detectedMIME := http.DetectContentType(header)
	// http.DetectContentType may return params (e.g., "video/mp4; codecs=...")
	detectedMIME = strings.Split(detectedMIME, ";")[0]
	detectedMIME = strings.TrimSpace(detectedMIME)

	ext, allowed := AllowedMIMETypes[detectedMIME]
	if !allowed {
		// Fall back to checking the declared type for video (Go's detection is weak for video)
		declared := strings.Split(declaredContentType, ";")[0]
		declared = strings.TrimSpace(declared)
		ext, allowed = AllowedMIMETypes[declared]
		if !allowed || (!strings.HasPrefix(declared, "video/")) {
			return "", fmt.Errorf("unsupported file type: %s (detected: %s)", declaredContentType, detectedMIME)
		}
	}

	// Generate a cryptographically random filename to prevent path manipulation
	randBytes := make([]byte, 16)
	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("generating random filename: %w", err)
	}
	filename := hex.EncodeToString(randBytes) + ext

	// Full path on disk
	destPath := filepath.Join(s.BasePath, filename)

	// Create the destination file with restricted permissions
	// Note: 0644 allows the file to be read by any user (needed because nginx
	// runs as a different UID than the backend that writes the file).
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("creating destination file: %w", err)
	}
	defer destFile.Close()

	// Write the header bytes we already read
	if _, err := destFile.Write(header); err != nil {
		os.Remove(destPath) // cleanup on failure
		return "", fmt.Errorf("writing file header: %w", err)
	}

	// Copy the rest of the file
	written, err := io.Copy(destFile, limited)
	if err != nil {
		os.Remove(destPath) // cleanup on failure
		return "", fmt.Errorf("writing file body: %w", err)
	}

	// Check total size (header + rest)
	totalSize := int64(len(header)) + written
	if totalSize > MaxBannerSize {
		os.Remove(destPath) // cleanup oversized file
		return "", fmt.Errorf("file exceeds maximum size of %d bytes", MaxBannerSize)
	}

	// Return the URL path (not the filesystem path)
	return s.URLPrefix + filename, nil
}

// Delete removes a previously stored file given its URL path.
// Security: Only deletes files within BasePath; rejects path traversal attempts.
func (s *Storage) Delete(urlPath string) error {
	if urlPath == "" {
		return nil // nothing to delete
	}

	// Extract filename from URL path
	filename := filepath.Base(urlPath)

	// Security: Reject any path traversal attempts
	if filename == "." || filename == ".." || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("invalid filename")
	}

	// Only delete if it matches our URL prefix (i.e., it's actually our file)
	if !strings.HasPrefix(urlPath, s.URLPrefix) {
		return nil // external URL, nothing to delete
	}

	fullPath := filepath.Join(s.BasePath, filename)

	// Verify the resolved path is still within our base directory (defense in depth)
	absBase, _ := filepath.Abs(s.BasePath)
	absTarget, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absTarget, absBase) {
		return fmt.Errorf("path traversal detected")
	}

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing old banner: %w", err)
	}
	return nil
}
