-- Add banner_url column to users table for custom profile banners.
-- Supports images (JPEG, PNG, WebP, GIF) and short videos (MP4, WebM).
-- The URL points to a locally stored file served via /uploads/.
ALTER TABLE users ADD COLUMN IF NOT EXISTS banner_url VARCHAR(512) DEFAULT '';
