"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import api from "@/lib/api";
import type { ShowcaseRepo, AcademicTag } from "@/types/showcase";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { useToast } from "@/components/ui/toast";
import {
  ExternalLink,
  Trash2,
  Pencil,
  MessageSquare,
  FolderGit2,
  Webhook,
  Check,
} from "lucide-react";

const ACADEMIC_TAGS: AcademicTag[] = [
  "coursework",
  "thesis",
  "hackathon",
  "personal_research",
  "team_project",
];

const TAG_COLORS: Record<string, string> = {
  coursework: "bg-blue-50 text-blue-700 border-blue-200",
  thesis: "bg-purple-50 text-purple-700 border-purple-200",
  hackathon: "bg-orange-50 text-orange-700 border-orange-200",
  personal_research: "bg-green-50 text-green-700 border-green-200",
  team_project: "bg-teal-50 text-teal-700 border-teal-200",
};

export function ShowcaseGrid() {
  const queryClient = useQueryClient();
  const { toast } = useToast();
  const [editingTagId, setEditingTagId] = useState<string | null>(null);

  const {
    data: repos,
    isLoading,
    isError,
    refetch,
  } = useQuery<ShowcaseRepo[]>({
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
    mutationFn: async ({
      repo,
      newTag,
    }: {
      repo: ShowcaseRepo;
      newTag: AcademicTag;
    }) => {
      await api.delete(`/showcase/${repo.id}`);
      try {
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
      } catch (postError) {
        try {
          await api.post("/showcase", {
            selections: [
              {
                repo_id: repo.github_repo_id,
                repo_name: repo.repo_name,
                full_name: repo.repo_full_name,
                tag: repo.academic_tag,
              },
            ],
          });
        } catch {
          // Rollback also failed
        }
        throw postError;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["showcaseRepos"] });
      toast("Tag updated successfully", "success");
      setEditingTagId(null);
    },
    onError: () => {
      queryClient.invalidateQueries({ queryKey: ["showcaseRepos"] });
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
          <Skeleton key={i} className="h-40 w-full rounded-2xl" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <Card className="border-red-100">
        <CardContent className="p-8 text-center">
          <div className="h-12 w-12 rounded-full bg-red-50 flex items-center justify-center mx-auto mb-3">
            <FolderGit2 className="h-6 w-6 text-red-400" />
          </div>
          <p className="text-gray-600 mb-3">
            Failed to load showcase repositories.
          </p>
          <Button variant="outline" size="sm" onClick={() => refetch()}>
            Try Again
          </Button>
        </CardContent>
      </Card>
    );
  }

  if (!repos || repos.length === 0) {
    return (
      <Card className="border-dashed border-2 border-gray-200">
        <CardContent className="p-10 text-center">
          <div className="h-14 w-14 rounded-full bg-gray-50 flex items-center justify-center mx-auto mb-4">
            <FolderGit2 className="h-7 w-7 text-gray-300" />
          </div>
          <h3 className="text-base font-medium text-gray-900 mb-1">
            No repositories yet
          </h3>
          <p className="text-sm text-gray-500 max-w-sm mx-auto">
            Add your best GitHub repositories below to start building your
            showcase portfolio.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2">
      {repos.map((repo) => (
        <Card
          key={repo.id}
          className="group hover:border-primary-200 hover:shadow-md transition-all duration-200 overflow-hidden"
        >
          {/* Colored top accent */}
          <div className="h-1 gradient-primary" />
          <CardContent className="p-5">
            <div className="flex items-start justify-between mb-3">
              <Link
                href={`/repos/${repo.id}/discussions`}
                className="font-semibold text-primary-600 hover:text-primary-700 hover:underline flex items-center gap-1.5 transition-colors"
              >
                <FolderGit2 className="h-4 w-4" />
                {repo.repo_name}
              </Link>
              <button
                onClick={() => handleRemove(repo)}
                disabled={deleteMutation.isPending}
                className="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition-all opacity-0 group-hover:opacity-100"
                title="Remove from showcase"
              >
                <Trash2 className="h-3.5 w-3.5" />
              </button>
            </div>

            <p className="text-sm text-gray-500 mb-3 line-clamp-2">
              {repo.description || "No description"}
            </p>

            <div className="flex items-center gap-2 flex-wrap mb-3">
              {repo.language && (
                <Badge variant="secondary" className="text-[10px]">
                  {repo.language}
                </Badge>
              )}
              {editingTagId === repo.id ? (
                <div className="flex flex-wrap gap-1">
                  {ACADEMIC_TAGS.map((tag) => (
                    <button
                      key={tag}
                      onClick={() => handleTagChange(repo, tag)}
                      className={`px-2 py-0.5 text-[10px] rounded-full border transition-all ${
                        tag === repo.academic_tag
                          ? "gradient-primary text-white border-transparent"
                          : "bg-white text-gray-600 border-gray-200 hover:border-primary-300"
                      }`}
                      disabled={updateTagMutation.isPending}
                    >
                      {tag.replace("_", " ")}
                    </button>
                  ))}
                </div>
              ) : (
                <div className="flex items-center gap-1.5">
                  <Badge
                    className={`text-[10px] ${TAG_COLORS[repo.academic_tag] || "bg-gray-50 text-gray-700"}`}
                  >
                    {repo.academic_tag.replace("_", " ")}
                  </Badge>
                  <button
                    onClick={() => setEditingTagId(repo.id)}
                    className="p-0.5 rounded text-gray-300 hover:text-primary-600 transition-colors"
                    title="Edit tag"
                  >
                    <Pencil className="h-3 w-3" />
                  </button>
                </div>
              )}
            </div>

            {/* Footer actions */}
            <div className="flex items-center gap-3 pt-3 border-t border-gray-50">
              <Link
                href={`/repos/${repo.id}/discussions`}
                className="text-xs text-gray-400 hover:text-primary-600 flex items-center gap-1 transition-colors"
              >
                <MessageSquare className="h-3 w-3" />
                Discussions
              </Link>
              <a
                href={repo.html_url || `https://github.com/${repo.repo_full_name}`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-xs text-gray-400 hover:text-primary-600 flex items-center gap-1 transition-colors"
              >
                <ExternalLink className="h-3 w-3" />
                GitHub
              </a>
              <div className="flex items-center gap-1 text-xs text-gray-400 ml-auto">
                <Webhook className="h-3 w-3" />
                <span>Webhook active</span>
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
