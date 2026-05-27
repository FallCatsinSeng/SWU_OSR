"use client";

import { useInfiniteQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import type { FeedResponse } from "@/types/activity";
import { ActivityCard } from "./ActivityCard";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { RefreshCw, Inbox, ArrowRight } from "lucide-react";
import Link from "next/link";

/**
 * Public activity feed — displays the global community feed
 * without requiring authentication. No sync button, no auth-dependent features.
 */
export function PublicActivityFeed() {
  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    isError,
    refetch,
  } = useInfiniteQuery<FeedResponse>({
    queryKey: ["publicActivityFeed"],
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

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-20 w-full rounded-geist-md" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <div className="h-12 w-12 rounded-geist-full bg-geist-error-soft dark:bg-neutral-800 flex items-center justify-center mx-auto mb-3">
            <RefreshCw className="h-5 w-5 text-geist-error dark:text-white" />
          </div>
          <p className="text-body-sm text-geist-body dark:text-white mb-4">
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
      <Card className="border border-dashed border-geist-hairline dark:border-neutral-700">
        <CardContent className="p-10 text-center">
          <div className="h-12 w-12 rounded-geist-full bg-geist-canvas-soft-2 dark:bg-neutral-800 flex items-center justify-center mx-auto mb-4">
            <Inbox className="h-6 w-6 text-geist-mute dark:text-white" />
          </div>
          <h3 className="text-body-md-strong text-geist-ink dark:text-white mb-1">
            No activity yet.
          </h3>
          <p className="text-body-sm text-geist-body dark:text-white max-w-sm mx-auto mb-6">
            The community feed will populate as members contribute to open source projects.
          </p>
          <Link href="/login">
            <Button variant="default" size="sm">
              Join the community
              <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
            </Button>
          </Link>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-3">
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
