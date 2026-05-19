"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import api from "@/lib/api";
import type { PublicProfile as PublicProfileType, AcademicIdentity } from "@/types/user";
import { useCurrentUser } from "@/hooks/useAuth";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogContent,
  DialogClose,
} from "@/components/ui/dialog";
import { ExternalLink, MessageSquare } from "lucide-react";

interface PublicProfileProps {
  alias: string;
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", { month: "long", year: "numeric" });
}

export function PublicProfile({ alias }: PublicProfileProps) {
  const { data: currentUser } = useCurrentUser();
  const [identityOpen, setIdentityOpen] = useState(false);
  const [identity, setIdentity] = useState<AcademicIdentity | null>(null);
  const [identityLoading, setIdentityLoading] = useState(false);

  const { data: profile, isLoading, isError } = useQuery<PublicProfileType>({
    queryKey: ["profile", alias],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: PublicProfileType }>(
        `/profiles/${alias}`
      );
      return data.data;
    },
  });

  const handleViewIdentity = async () => {
    setIdentityLoading(true);
    try {
      const { data } = await api.get<{ ok: boolean; data: AcademicIdentity }>(
        `/profiles/${alias}/identity`
      );
      setIdentity(data.data);
      setIdentityOpen(true);
    } catch {
      setIdentity(null);
      setIdentityOpen(true);
    } finally {
      setIdentityLoading(false);
    }
  };

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
              <p className="text-sm text-gray-400 mt-1">
                Joined {formatDate(profile.created_at)}
              </p>
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
              {currentUser && (
                <div className="mt-3">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleViewIdentity}
                    disabled={identityLoading}
                  >
                    {identityLoading ? "Loading..." : "View Real Identity"}
                  </Button>
                </div>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {profile.stats && (
        <div>
          <h2 className="text-xl font-semibold text-gray-900 mb-3">Stats</h2>
          <div className="grid grid-cols-2 sm:grid-cols-5 gap-3">
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-2xl font-bold text-gray-900">
                  {profile.stats.total_commits}
                </p>
                <p className="text-xs text-gray-500">Commits</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-2xl font-bold text-gray-900">
                  {profile.stats.total_repos}
                </p>
                <p className="text-xs text-gray-500">Repos</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-2xl font-bold text-gray-900">
                  {profile.stats.languages.length}
                </p>
                <p className="text-xs text-gray-500">Languages</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-2xl font-bold text-gray-900">
                  {profile.stats.active_days}
                </p>
                <p className="text-xs text-gray-500">Active Days</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <p className="text-2xl font-bold text-gray-900">
                  {profile.stats.current_streak}
                </p>
                <p className="text-xs text-gray-500">Streak</p>
              </CardContent>
            </Card>
          </div>
        </div>
      )}

      <div>
        <h2 className="text-xl font-semibold text-gray-900 mb-3">Badges</h2>
        <p className="text-sm text-gray-400">No badges yet</p>
      </div>

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
                  <div className="mt-2">
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
        </div>
      )}

      <Dialog open={identityOpen} onOpenChange={setIdentityOpen}>
        <DialogClose onClose={() => setIdentityOpen(false)} />
        <DialogHeader>
          <DialogTitle>Real Identity</DialogTitle>
        </DialogHeader>
        <DialogContent>
          {identity ? (
            <div className="space-y-2">
              <p>
                <span className="font-medium text-gray-700">Full Name:</span>{" "}
                {identity.full_name}
              </p>
              <p>
                <span className="font-medium text-gray-700">NIM:</span>{" "}
                {identity.nim}
              </p>
              <p>
                <span className="font-medium text-gray-700">Major:</span>{" "}
                {identity.major}
              </p>
              <p>
                <span className="font-medium text-gray-700">Semester:</span>{" "}
                {identity.semester}
              </p>
            </div>
          ) : (
            <p className="text-gray-500">Unable to load identity information.</p>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
