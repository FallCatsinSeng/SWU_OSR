import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { resolveUploadUrl } from '../url';

describe('resolveUploadUrl', () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  it('returns empty string for falsy input', () => {
    expect(resolveUploadUrl(undefined)).toBe('');
    expect(resolveUploadUrl(null)).toBe('');
    expect(resolveUploadUrl('')).toBe('');
  });

  it('returns full HTTPS URLs as-is', () => {
    const url = 'https://cdn.example.com/uploads/banner.jpg';
    expect(resolveUploadUrl(url)).toBe(url);
  });

  it('returns full HTTP URLs as-is', () => {
    const url = 'http://localhost:8080/uploads/banner.jpg';
    expect(resolveUploadUrl(url)).toBe(url);
  });

  it('returns relative path as-is when no NEXT_PUBLIC_UPLOADS_URL is set', () => {
    delete process.env.NEXT_PUBLIC_UPLOADS_URL;
    const path = '/uploads/banners/abc123.jpg';
    expect(resolveUploadUrl(path)).toBe(path);
  });

  it('prepends NEXT_PUBLIC_UPLOADS_URL when set', () => {
    process.env.NEXT_PUBLIC_UPLOADS_URL = 'https://cdn.swu-osr.ac.id';
    const path = '/uploads/banners/abc123.jpg';
    const result = resolveUploadUrl(path);
    expect(result).toBe('https://cdn.swu-osr.ac.id/uploads/banners/abc123.jpg');
  });

  it('handles NEXT_PUBLIC_UPLOADS_URL with trailing slash', () => {
    process.env.NEXT_PUBLIC_UPLOADS_URL = 'https://cdn.swu-osr.ac.id/';
    const path = '/uploads/banners/abc123.jpg';
    const result = resolveUploadUrl(path);
    // Should not double-slash
    expect(result).not.toContain('//uploads');
    expect(result).toBe('https://cdn.swu-osr.ac.id/uploads/banners/abc123.jpg');
  });
});
