"use client";

import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";
import type { PublicProfile as PublicProfileType } from "@/types/user";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { ExternalLink } from "lucide-react";

interface PublicProfileProps {
  alias: string;
}

export function PublicProfile({ alias }: PublicProfileProps) {
  const { data: profile, isLoading, isError } = useQuery<PublicProfileType>({
    queryKey: ["profile", alias],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: PublicProfileType }>(
        `/profiles/${alias}`
      );
      return data.data;
    },
  });

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-32 w-full" />
        <Skeleton className="h-48 w-full" />
      </div>
    );
  }

  if (isError || !profile) {
    return <p className="text-center text-gray-500">Profile not found.</p>;
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardContent className="p-6">
          <div className="flex items-start gap-4">
            <Avatar
              src={profile.avatar_url}
              alt={profile.alias}
              fallback={profile.alias.charAt(0).toUpperCase()}
              size="lg"
            />
            <div className="flex-1">
              <h1 className="text-2xl font-bold text-gray-900">
                {profile.alias}
              </h1>
              {profile.bio && (
                <p className="text-gray-600 mt-1">{profile.bio}</p>
              )}
              <div className="flex items-center gap-2 mt-2">
                <Badge>{profile.role}</Badge>
                {profile.github_username && (
                  <a
                    href={`https://github.com/${profile.github_username}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-blue-600 hover:underline flex items-center gap-1"
                  >
                    @{profile.github_username}
                    <ExternalLink className="h-3 w-3" />
                  </a>
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {profile.showcase_repos && profile.showcase_repos.length > 0 && (
        <div>
          <h2 className="text-xl font-semibold text-gray-900 mb-3">
            Showcase Repositories
          </h2>
          <div className="grid gap-4 sm:grid-cols-2">
            {profile.showcase_repos.map((repo) => (
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
                  <p className="text-sm text-gray-600 mb-2">
                    {repo.description || "No description"}
                  </p>
                  <div className="flex items-center gap-2">
                    {repo.language && (
                      <Badge variant="secondary">{repo.language}</Badge>
                    )}
                    <Badge>{repo.academic_tag}</Badge>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
