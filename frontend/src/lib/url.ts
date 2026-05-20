/**
 * Resolves an upload path (e.g., "/uploads/banners/abc.jpg") to a full URL
 * that points to the backend server where the file is actually served.
 *
 * In production (behind nginx), relative paths work fine because nginx serves
 * /uploads/* directly. In development (accessing frontend at localhost:3000),
 * we need to resolve against the backend origin.
 */
export function resolveUploadUrl(path: string | undefined | null): string {
  if (!path) return "";

  // Already a full URL (e.g., https://...)
  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }

  // In production behind nginx, NEXT_PUBLIC_API_URL is usually "/api" or empty,
  // so relative paths work. In development, it's "http://localhost:8080".
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "";

  // If API URL is a full origin (starts with http), use its origin for uploads
  if (apiUrl.startsWith("http")) {
    try {
      const origin = new URL(apiUrl).origin;
      return `${origin}${path}`;
    } catch {
      return path;
    }
  }

  // Otherwise (behind proxy/nginx), relative path is fine
  return path;
}
