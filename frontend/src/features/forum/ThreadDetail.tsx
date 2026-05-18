"use client";

import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import type { Thread, Comment } from "@/types/forum";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { CreateCommentForm } from "./CreateCommentForm";

interface ThreadDetailProps {
  repoId: string;
  threadId: string;
}

export function ThreadDetail({ repoId, threadId }: ThreadDetailProps) {
  const { data, isLoading } = useQuery<{ thread: Thread; comments: Comment[] }>({
    queryKey: ["thread", threadId],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: { thread: Thread; comments: Comment[] } }>(
        `/threads/${threadId}`
      );
      return data.data;
    },
  });

  const thread = data?.thread;
  const comments = data?.comments;

  if (isLoading) {
    return <Skeleton className="h-64 w-full" />;
  }

  if (!thread) {
    return <p className="text-center text-gray-500">Thread not found.</p>;
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>{thread.title}</CardTitle>
          <p className="text-sm text-gray-400">
            {new Date(thread.created_at).toLocaleDateString()}
          </p>
        </CardHeader>
        <CardContent>
          <p className="text-gray-700 whitespace-pre-wrap">{thread.body}</p>
        </CardContent>
      </Card>

      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          Comments ({thread.comment_count})
        </h3>
        {comments && comments.length > 0 ? (
          <div className="space-y-3">
            {comments.map((comment) => (
              <Card key={comment.id}>
                <CardContent className="p-4">
                  <p className="text-sm text-gray-700 whitespace-pre-wrap">
                    {comment.body}
                  </p>
                  <p className="text-xs text-gray-400 mt-2">
                    {new Date(comment.created_at).toLocaleDateString()}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : (
          <p className="text-gray-500 text-sm">No comments yet.</p>
        )}
      </div>

      <CreateCommentForm repoId={repoId} threadId={threadId} />
    </div>
  );
}
