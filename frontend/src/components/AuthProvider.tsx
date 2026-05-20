"use client";

import * as React from "react";
import { rehydrateToken, getAccessToken } from "@/lib/auth";

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
 * AuthProvider rehydrates the access token from the httpOnly refresh_token
 * cookie on mount. Until rehydration completes, children see isReady=false.
 *
 * It also exposes setAuthenticated so that login mutations can update
 * the context without requiring a page reload.
 */
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isReady, setIsReady] = React.useState(false);
  const [isAuthenticated, setIsAuthenticated] = React.useState(false);

  React.useEffect(() => {
    let mounted = true;

    async function init() {
      // Attempt to rehydrate token from refresh cookie
      const token = await rehydrateToken();
      if (mounted) {
        setIsAuthenticated(!!token);
        setIsReady(true);
      }
    }

    // If there's already a token in memory (e.g., just logged in), skip rehydration
    if (getAccessToken()) {
      setIsAuthenticated(true);
      setIsReady(true);
    } else {
      init();
    }

    return () => {
      mounted = false;
    };
  }, []);

  const value = React.useMemo(
    () => ({ isReady, isAuthenticated, setAuthenticated: setIsAuthenticated }),
    [isReady, isAuthenticated]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
