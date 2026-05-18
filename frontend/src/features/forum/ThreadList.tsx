"use client";

import { useInfiniteQuery } from "@tanstack/react-query";
import Link from "next/link";
import api from "@/lib/api";
import type { ThreadList as ThreadListType } from "@/types/forum";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { MessageSquare } from "lucide-react";

interface ThreadListProps {
  repoId: string;
}

export function ThreadList({ repoId }: ThreadListProps) {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading, isError } =
    useInfiniteQuery<ThreadListType>({
      queryKey: ["threads", repoId],
      queryFn: async ({ pageParam }) => {
        const params = pageParam ? { cursor: pageParam, limit: 20 } : { limit: 20 };
        const { data } = await api.get<{ ok: boolean; data: ThreadListType }>(
          `/repos/${repoId}/threads`,
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
          <Skeleton key={i} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  if (isError) {
    return <p className="text-center text-gray-500">Failed to load threads.</p>;
  }

  const threads = data?.pages.flatMap((page) => page.threads) ?? [];

  if (threads.length === 0) {
    return (
      <p className="text-center text-gray-500">
        No discussions yet. Start one above!
      </p>
    );
  }

  return (
    <div className="space-y-3">
      {threads.map((thread) => (
        <Link
          key={thread.id}
          href={`/repos/${repoId}/discussions/${thread.id}`}
        >
          <Card className="hover:border-blue-200 transition-colors cursor-pointer">
            <CardContent className="p-4">
              <div className="flex items-start justify-between">
                <div>
                  <h3 className="font-medium text-gray-900">{thread.title}</h3>
                  <p className="text-sm text-gray-500 mt-1 line-clamp-1">
                    {thread.body}
                  </p>
                </div>
                <Badge variant="secondary" className="flex items-center gap-1">
                  <MessageSquare className="h-3 w-3" />
                  {thread.comment_count}
                </Badge>
              </div>
            </CardContent>
          </Card>
        </Link>
      ))}
      {hasNextPage && (
        <div className="flex justify-center pt-4">
          <Button
            variant="outline"
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
          >
            {isFetchingNextPage ? "Loading..." : "Load More"}
          </Button>
        </div>
      )}
    </div>
  );
}
