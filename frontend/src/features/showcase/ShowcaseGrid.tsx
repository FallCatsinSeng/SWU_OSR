"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import api from "@/lib/api";
import type { ShowcaseRepo, AcademicTag } from "@/types/showcase";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { useToast } from "@/components/ui/toast";
import { ExternalLink, Trash2, Pencil, MessageSquare } from "lucide-react";

const ACADEMIC_TAGS: AcademicTag[] = [
  "coursework",
  "thesis",
  "hackathon",
  "personal_research",
  "team_project",
];

export function ShowcaseGrid() {
  const queryClient = useQueryClient();
  const { toast } = useToast();
  const [editingTagId, setEditingTagId] = useState<string | null>(null);

  const { data: repos, isLoading, isError, refetch } = useQuery<ShowcaseRepo[]>({
    queryKey: ["showcaseRepos"],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: ShowcaseRepo[] }>(
        "/showcase"
      );
      return data.data;
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/showcase/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["showcaseRepos"] });
      toast("Repository removed from showcase", "success");
    },
    onError: () => {
      toast("Failed to remove repository", "error");
    },
  });

  const updateTagMutation = useMutation({
    mutationFn: async ({ repo, newTag }: { repo: ShowcaseRepo; newTag: AcademicTag }) => {
      await api.delete(`/showcase/${repo.id}`);
      await api.post("/showcase", {
        selections: [
          {
            repo_id: repo.github_repo_id,
            repo_name: repo.repo_name,
            full_name: repo.repo_full_name,
            tag: newTag,
          },
        ],
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["showcaseRepos"] });
      toast("Tag updated successfully", "success");
      setEditingTagId(null);
    },
    onError: () => {
      toast("Failed to update tag", "error");
      setEditingTagId(null);
    },
  });

  const handleRemove = (repo: ShowcaseRepo) => {
    if (window.confirm(`Remove "${repo.repo_name}" from your showcase?`)) {
      deleteMutation.mutate(repo.id);
    }
  };

  const handleTagChange = (repo: ShowcaseRepo, newTag: AcademicTag) => {
    if (newTag !== repo.academic_tag) {
      updateTagMutation.mutate({ repo, newTag });
    } else {
      setEditingTagId(null);
    }
  };

  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-32 w-full" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardContent className="p-6 text-center">
          <p className="text-gray-600 mb-3">Failed to load showcase repositories.</p>
          <Button variant="outline" size="sm" onClick={() => refetch()}>
            Try Again
          </Button>
        </CardContent>
      </Card>
    );
  }

  if (!repos || repos.length === 0) {
    return (
      <p className="text-center text-gray-500">
        You have not added any repositories to your showcase yet.
      </p>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2">
      {repos.map((repo) => (
        <Card key={repo.id}>
          <CardHeader className="pb-2">
            <CardTitle className="text-base flex items-center justify-between">
              <a
                href={repo.html_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:underline flex items-center gap-1"
              >
                {repo.repo_name}
                <ExternalLink className="h-3 w-3" />
              </a>
              <Button
                variant="destructive"
                size="sm"
                onClick={() => handleRemove(repo)}
                disabled={deleteMutation.isPending}
                className="h-7 px-2"
              >
                <Trash2 className="h-3 w-3" />
              </Button>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-3">
              {repo.description || "No description"}
            </p>
            <div className="flex items-center gap-2 flex-wrap">
              {repo.language && (
                <Badge variant="secondary">{repo.language}</Badge>
              )}
              {editingTagId === repo.id ? (
                <div className="flex flex-wrap gap-1">
                  {ACADEMIC_TAGS.map((tag) => (
                    <button
                      key={tag}
                      onClick={() => handleTagChange(repo, tag)}
                      className={`px-2 py-0.5 text-xs rounded-full border transition-colors ${
                        tag === repo.academic_tag
                          ? "bg-blue-600 text-white border-blue-600"
                          : "bg-white text-gray-600 border-gray-300 hover:border-blue-300"
                      }`}
                      disabled={updateTagMutation.isPending}
                    >
                      {tag.replace("_", " ")}
                    </button>
                  ))}
                </div>
              ) : (
                <div className="flex items-center gap-1">
                  <Badge>{repo.academic_tag.replace("_", " ")}</Badge>
                  <button
                    onClick={() => setEditingTagId(repo.id)}
                    className="text-gray-400 hover:text-gray-600"
                    title="Edit tag"
                  >
                    <Pencil className="h-3 w-3" />
                  </button>
                </div>
              )}
            </div>
            <div className="mt-3">
              <Link
                href={`/repos/${repo.id}/discussions`}
                className="text-sm text-gray-500 hover:text-blue-600 flex items-center gap-1"
              >
                <MessageSquare className="h-3 w-3" />
                Discussions
              </Link>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
