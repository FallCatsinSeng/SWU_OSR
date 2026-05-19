"use client";

import { useEffect, useRef } from "react";
import Link from "next/link";
import {
  useInfiniteQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import api from "@/lib/api";
import type { FeedResponse } from "@/types/activity";
import { useCurrentUser } from "@/hooks/useAuth";
import { ActivityCard } from "./ActivityCard";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useToast } from "@/components/ui/toast";
import { FolderGit2, RefreshCw, Inbox } from "lucide-react";

export function ActivityFeed() {
  const { data: user } = useCurrentUser();
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    isError,
    refetch,
  } = useInfiniteQuery<FeedResponse>({
    queryKey: ["activityFeed"],
    queryFn: async ({ pageParam }) => {
      const params = pageParam
        ? { cursor: pageParam, limit: 20 }
        : { limit: 20 };
      const { data } = await api.get<{ ok: boolean; data: FeedResponse }>(
        "/feed",
        { params }
      );
      return data.data;
    },
    initialPageParam: "",
    getNextPageParam: (lastPage) =>
      lastPage.has_more ? lastPage.next_cursor : undefined,
  });

  const syncMutation = useMutation({
    mutationFn: async () => {
      const { data } = await api.post<{
        ok: boolean;
        data: { synced: number };
      }>("/activity/sync");
      return data.data;
    },
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: ["activityFeed"] });
      if (result.synced > 0) {
        toast(`Synced ${result.synced} new activities from GitHub`, "success");
      } else {
        toast("Already up to date — no new activities found", "success");
      }
    },
    onError: () => {
      toast("Failed to sync activity from GitHub", "error");
    },
  });

  // Auto-sync on first load when feed is empty
  const autoSyncRef = useRef(false);
  useEffect(() => {
    if (autoSyncRef.current) return;
    if (isLoading || isError) return;
    const items = data?.pages.flatMap((page) => page.items) ?? [];
    if (items.length === 0 && user && !syncMutation.isPending) {
      autoSyncRef.current = true;
      syncMutation.mutate();
    }
  }, [data, user, isLoading, isError]);

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-20 w-full rounded-geist-md" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <div className="h-12 w-12 rounded-geist-full bg-geist-error-soft flex items-center justify-center mx-auto mb-3">
            <RefreshCw className="h-5 w-5 text-geist-error" />
          </div>
          <p className="text-body-sm text-geist-body mb-4">
            Failed to load activity feed.
          </p>
          <Button variant="outline" size="sm" onClick={() => refetch()}>
            Try again
          </Button>
        </CardContent>
      </Card>
    );
  }

  const items = data?.pages.flatMap((page) => page.items) ?? [];

  if (items.length === 0) {
    return (
      <Card className="border border-dashed border-geist-hairline">
        <CardContent className="p-10 text-center">
          <div className="h-12 w-12 rounded-geist-full bg-geist-canvas-soft-2 flex items-center justify-center mx-auto mb-4">
            <Inbox className="h-6 w-6 text-geist-mute" />
          </div>
          <h3 className="text-body-md-strong text-geist-ink mb-1">
            No activity yet.
          </h3>
          <p className="text-body-sm text-geist-body max-w-sm mx-auto mb-6">
            {user
              ? "Sync your GitHub activity or add repos to your showcase to start tracking contributions."
              : "Sign in and add repos to your showcase to start tracking open source contributions."}
          </p>
          <div className="flex items-center justify-center gap-3">
            {user && (
              <Button
                variant="default"
                size="sm"
                onClick={() => syncMutation.mutate()}
                disabled={syncMutation.isPending}
              >
                {syncMutation.isPending ? (
                  <RefreshCw className="mr-1.5 h-3.5 w-3.5 animate-spin" />
                ) : (
                  <RefreshCw className="mr-1.5 h-3.5 w-3.5" />
                )}
                {syncMutation.isPending
                  ? "Syncing..."
                  : "Sync GitHub activity"}
              </Button>
            )}
            <Link href={user ? "/showcase" : "/login"}>
              <Button variant="outline" size="sm">
                <FolderGit2 className="mr-1.5 h-3.5 w-3.5" />
                {user ? "Go to Showcase" : "Sign In"}
              </Button>
            </Link>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-3">
      {/* Sync button */}
      {user && (
        <div className="flex justify-end mb-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => syncMutation.mutate()}
            disabled={syncMutation.isPending}
          >
            {syncMutation.isPending ? (
              <RefreshCw className="mr-1.5 h-3 w-3 animate-spin" />
            ) : (
              <RefreshCw className="mr-1.5 h-3 w-3" />
            )}
            {syncMutation.isPending ? "Syncing..." : "Sync from GitHub"}
          </Button>
        </div>
      )}

      {items.map((item) => (
        <ActivityCard key={item.id} item={item} />
      ))}

      {hasNextPage && (
        <div className="flex justify-center pt-4">
          <Button
            variant="outline"
            size="sm"
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
          >
            {isFetchingNextPage ? (
              <>
                <RefreshCw className="mr-1.5 h-3.5 w-3.5 animate-spin" />
                Loading...
              </>
            ) : (
              "Load more"
            )}
          </Button>
        </div>
      )}
    </div>
  );
}
