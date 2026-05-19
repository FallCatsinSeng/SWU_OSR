"use client";

import Link from "next/link";
import { useInfiniteQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import type { FeedResponse } from "@/types/activity";
import { useCurrentUser } from "@/hooks/useAuth";
import { ActivityCard } from "./ActivityCard";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useToast } from "@/components/ui/toast";
import { FolderGit2, RefreshCw, Inbox, Zap } from "lucide-react";

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
      const { data } = await api.post<{ ok: boolean; data: { synced: number } }>(
        "/activity/sync"
      );
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

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-20 w-full rounded-xl" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <Card className="border-red-100">
        <CardContent className="p-8 text-center">
          <div className="h-12 w-12 rounded-full bg-red-50 flex items-center justify-center mx-auto mb-3">
            <RefreshCw className="h-6 w-6 text-red-400" />
          </div>
          <p className="text-gray-600 mb-3">Failed to load activity feed.</p>
          <Button variant="outline" size="sm" onClick={() => refetch()}>
            Try Again
          </Button>
        </CardContent>
      </Card>
    );
  }

  const items = data?.pages.flatMap((page) => page.items) ?? [];

  if (items.length === 0) {
    return (
      <Card className="border-dashed border-2 border-gray-200">
        <CardContent className="p-10 text-center">
          <div className="h-14 w-14 rounded-full bg-gray-50 flex items-center justify-center mx-auto mb-4">
            <Inbox className="h-7 w-7 text-gray-300" />
          </div>
          <h3 className="text-base font-medium text-gray-900 mb-1">
            No activity yet
          </h3>
          <p className="text-sm text-gray-500 max-w-sm mx-auto mb-4">
            {user
              ? "Sync your GitHub activity or add repos to your showcase to start tracking contributions."
              : "Sign in and add repos to your showcase to start tracking open source contributions."}
          </p>
          <div className="flex items-center justify-center gap-3">
            {user && (
              <Button
                variant="default"
                size="sm"
                className="gap-1.5 gradient-primary text-white border-0"
                onClick={() => syncMutation.mutate()}
                disabled={syncMutation.isPending}
              >
                {syncMutation.isPending ? (
                  <RefreshCw className="h-3.5 w-3.5 animate-spin" />
                ) : (
                  <Zap className="h-3.5 w-3.5" />
                )}
                {syncMutation.isPending ? "Syncing..." : "Sync GitHub Activity"}
              </Button>
            )}
            <Link href={user ? "/showcase" : "/login"}>
              <Button variant="outline" size="sm" className="gap-1.5">
                <FolderGit2 className="h-3.5 w-3.5" />
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
      {/* Sync button at top when there are items */}
      {user && (
        <div className="flex justify-end mb-2">
          <Button
            variant="outline"
            size="sm"
            className="gap-1.5 text-xs"
            onClick={() => syncMutation.mutate()}
            disabled={syncMutation.isPending}
          >
            {syncMutation.isPending ? (
              <RefreshCw className="h-3 w-3 animate-spin" />
            ) : (
              <RefreshCw className="h-3 w-3" />
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
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
            className="gap-1.5"
          >
            {isFetchingNextPage ? (
              <>
                <RefreshCw className="h-3.5 w-3.5 animate-spin" />
                Loading...
              </>
            ) : (
              "Load More"
            )}
          </Button>
        </div>
      )}
    </div>
  );
}
