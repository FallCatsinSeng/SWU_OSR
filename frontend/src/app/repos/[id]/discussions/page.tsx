"use client";

import { ThreadList } from "@/features/forum/ThreadList";
import { CreateThreadForm } from "@/features/forum/CreateThreadForm";

interface DiscussionsPageProps {
  params: { id: string };
}

export default function DiscussionsPage({ params }: DiscussionsPageProps) {
  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Discussions</h1>
      </div>
      <CreateThreadForm repoId={params.id} />
      <div className="mt-6">
        <ThreadList repoId={params.id} />
      </div>
    </div>
  );
}
