"use client";

import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import type { ShowcaseRepo } from "@/types/showcase";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { ExternalLink } from "lucide-react";

export function ShowcaseGrid() {
  const { data: repos, isLoading } = useQuery<ShowcaseRepo[]>({
    queryKey: ["showcaseRepos"],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: ShowcaseRepo[] }>(
        "/showcase/my-repos"
      );
      return data.data;
    },
  });

  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-32 w-full" />
        ))}
      </div>
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
            <CardTitle className="text-base">
              <a
                href={repo.html_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:underline flex items-center gap-1"
              >
                {repo.repo_name}
                <ExternalLink className="h-3 w-3" />
              </a>
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
              <Badge>{repo.academic_tag.replace("_", " ")}</Badge>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
