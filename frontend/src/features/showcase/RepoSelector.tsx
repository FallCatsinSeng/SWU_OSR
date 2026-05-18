"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import type { Repository, AcademicTag, ShowcaseSelection } from "@/types/showcase";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { useToast } from "@/components/ui/toast";
import { Check } from "lucide-react";

const ACADEMIC_TAGS: AcademicTag[] = [
  "coursework",
  "thesis",
  "hackathon",
  "personal_research",
  "team_project",
];

export function RepoSelector() {
  const [selections, setSelections] = useState<Map<number, AcademicTag>>(new Map());
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const { data: repos, isLoading } = useQuery<Repository[]>({
    queryKey: ["availableRepos"],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: Repository[] }>(
        "/showcase/available-repos"
      );
      return data.data;
    },
  });

  const saveShowcase = useMutation({
    mutationFn: async (items: ShowcaseSelection[]) => {
      const { data } = await api.post("/showcase/select", { repos: items });
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["showcaseRepos"] });
      toast("Showcase updated successfully", "success");
      setSelections(new Map());
    },
    onError: () => {
      toast("Failed to update showcase", "error");
    },
  });

  const toggleRepo = (repo: Repository, tag: AcademicTag) => {
    const next = new Map(selections);
    if (next.has(repo.id) && next.get(repo.id) === tag) {
      next.delete(repo.id);
    } else {
      next.set(repo.id, tag);
    }
    setSelections(next);
  };

  const handleSave = () => {
    if (selections.size === 0) return;
    const items: ShowcaseSelection[] = [];
    if (repos) {
      for (const repo of repos) {
        const tag = selections.get(repo.id);
        if (tag) {
          items.push({
            repo_id: repo.id,
            repo_name: repo.name,
            full_name: repo.full_name,
            tag,
          });
        }
      }
    }
    saveShowcase.mutate(items);
  };

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-20 w-full" />
        ))}
      </div>
    );
  }

  if (!repos || repos.length === 0) {
    return (
      <p className="text-center text-gray-500">
        No repositories found. Make sure your GitHub account is linked.
      </p>
    );
  }

  return (
    <div className="space-y-4">
      <div className="space-y-3">
        {repos.map((repo) => (
          <Card
            key={repo.id}
            className={selections.has(repo.id) ? "border-blue-300 bg-blue-50" : ""}
          >
            <CardContent className="p-4">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h3 className="font-medium text-gray-900">{repo.name}</h3>
                  <p className="text-sm text-gray-500 mt-1">
                    {repo.description || "No description"}
                  </p>
                  {repo.language && (
                    <Badge variant="secondary" className="mt-2">
                      {repo.language}
                    </Badge>
                  )}
                </div>
                {selections.has(repo.id) && (
                  <Check className="h-5 w-5 text-blue-600" />
                )}
              </div>
              <div className="flex flex-wrap gap-2 mt-3">
                {ACADEMIC_TAGS.map((tag) => (
                  <button
                    key={tag}
                    onClick={() => toggleRepo(repo, tag)}
                    className={`px-2 py-1 text-xs rounded-full border transition-colors ${
                      selections.get(repo.id) === tag
                        ? "bg-blue-600 text-white border-blue-600"
                        : "bg-white text-gray-600 border-gray-300 hover:border-blue-300"
                    }`}
                  >
                    {tag.replace("_", " ")}
                  </button>
                ))}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
      {selections.size > 0 && (
        <div className="flex justify-end">
          <Button onClick={handleSave} disabled={saveShowcase.isPending}>
            {saveShowcase.isPending
              ? "Saving..."
              : `Add ${selections.size} to Showcase`}
          </Button>
        </div>
      )}
    </div>
  );
}
