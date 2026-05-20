// Security: Access token stored in memory only — never in localStorage.
// On page refresh the token is lost and transparently recovered via the
// httpOnly refresh_token cookie hitting /auth/refresh.

let accessToken: string | null = null;

const SESSION_ID_KEY = "swu_osr_session_id";

export function getAccessToken(): string | null {
  return accessToken;
}

export function setAccessToken(token: string): void {
  accessToken = token;
}

export function getSessionId(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(SESSION_ID_KEY);
}

export function setSessionId(sessionId: string): void {
  if (typeof window === "undefined") return;
  localStorage.setItem(SESSION_ID_KEY, sessionId);
}

export function clearTokens(): void {
  accessToken = null;
  if (typeof window === "undefined") return;
  localStorage.removeItem(SESSION_ID_KEY);
}

/**
 * Attempt to rehydrate the access token from the server using the
 * httpOnly refresh_token cookie. Call this once on application init.
 */
export async function rehydrateToken(): Promise<string | null> {
  try {
    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL || "/api"}/auth/refresh`,
      {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: "{}",
      }
    );
    if (!res.ok) return null;
    const data = await res.json();
    const token = data.data?.access_token ?? data.access_token;
    if (token) {
      accessToken = token;
      return token;
    }
    return null;
  } catch {
    return null;
  }
}
