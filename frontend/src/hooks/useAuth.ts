"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { setAccessToken, setSessionId, clearTokens } from "@/lib/auth";
import type { LoginInput, PendingSession, AuthResult } from "@/types/auth";
import type { User } from "@/types/user";

export function useSIAKADLogin() {
  return useMutation({
    mutationFn: async (input: LoginInput) => {
      const { data } = await api.post<{ ok: boolean; data: PendingSession }>(
        "/auth/siakad-login",
        input
      );
      return data.data;
    },
    onSuccess: (data) => {
      setSessionId(data.session_id);
      window.location.href = data.github_oauth_url;
    },
  });
}

export function useGitHubCallback() {
  return useMutation({
    mutationFn: async (params: { code: string; state: string }) => {
      const { data } = await api.post<{ ok: boolean; data: AuthResult }>(
        "/auth/github-callback",
        params
      );
      return data.data;
    },
    onSuccess: (data) => {
      setAccessToken(data.access_token);
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      await api.post("/auth/logout");
    },
    onSuccess: () => {
      clearTokens();
      queryClient.clear();
      window.location.href = "/";
    },
  });
}

export function useCurrentUser() {
  return useQuery<User>({
    queryKey: ["currentUser"],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: User }>("/auth/me");
      return data.data;
    },
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}
