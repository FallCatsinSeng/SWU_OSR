"use client";

import * as React from "react";
import { rehydrateToken, getAccessToken, hasLoggedInHint } from "@/lib/auth";

interface AuthContextValue {
  /** True once the initial token rehydration attempt has completed. */
  isReady: boolean;
  /** True when a valid access token exists in memory. */
  isAuthenticated: boolean;
  /** Call after a successful login to update context state. */
  setAuthenticated: (value: boolean) => void;
}

const AuthContext = React.createContext<AuthContextValue>({
  isReady: false,
  isAuthenticated: false,
  setAuthenticated: () => {},
});

export function useAuthContext() {
  return React.useContext(AuthContext);
}

/**
 * Compute the initial auth state synchronously so the very first render
 * already knows if the user is authenticated (from sessionStorage).
 * This eliminates the single-frame flash that occurs when state starts
 * as false and is corrected in useEffect.
 */
function getInitialAuthState(): { isReady: boolean; isAuthenticated: boolean } {
  // On the server (SSR), we can't access sessionStorage/localStorage
  if (typeof window === "undefined") {
    return { isReady: false, isAuthenticated: false };
  }
  // If we have a token in memory or sessionStorage, we're instantly ready
  const token = getAccessToken();
  if (token) {
    return { isReady: true, isAuthenticated: true };
  }
  // No token available synchronously — will need async rehydration
  return { isReady: false, isAuthenticated: false };
}

/**
 * AuthProvider rehydrates the access token on mount.
 *
 * Priority order:
 * 1. If token already exists in sessionStorage → instant ready on first render.
 * 2. Otherwise, check the httpOnly refresh cookie via /auth/refresh.
 *
 * While waiting for (2), the `hasLoggedInHint()` flag tells consuming
 * components whether the user was previously logged in, so they can show
 * a loading skeleton instead of flashing the unauthenticated landing page.
 */
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isReady, setIsReady] = React.useState(() => getInitialAuthState().isReady);
  const [isAuthenticated, setIsAuthenticated] = React.useState(() => getInitialAuthState().isAuthenticated);

  React.useEffect(() => {
    // If already resolved synchronously (token was in sessionStorage), skip
    if (isReady) return;

    let mounted = true;

    async function init() {
      const token = await rehydrateToken();
      if (mounted) {
        setIsAuthenticated(!!token);
        setIsReady(true);
      }
    }

    init();

    return () => {
      mounted = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const value = React.useMemo(
    () => ({ isReady, isAuthenticated, setAuthenticated: setIsAuthenticated }),
    [isReady, isAuthenticated]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
