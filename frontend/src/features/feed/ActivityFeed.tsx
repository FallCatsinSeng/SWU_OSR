"use client";

import Link from "next/link";
import { useInfiniteQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import type { FeedResponse } from "@/types/activity";
import { ActivityCard } from "./ActivityCard";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { FolderGit2, RefreshCw, Inbox } from "lucide-react";

export function ActivityFeed() {
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
            Add repos to your showcase to start tracking your open source
            contributions.
          </p>
          <Link href="/showcase">
            <Button variant="outline" size="sm" className="gap-1.5">
              <FolderGit2 className="h-3.5 w-3.5" />
              Go to Showcase
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
