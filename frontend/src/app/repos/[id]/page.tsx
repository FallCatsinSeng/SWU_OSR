"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import api from "@/lib/api";
import type { ShowcaseRepo } from "@/types/showcase";
import type { ActivityItem } from "@/types/activity";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import {
  ArrowLeft,
  ExternalLink,
  FolderGit2,
  MessageSquare,
  GitBranch,
  GitCommit,
  GitPullRequest,
  Tag,
  Code2,
  Clock,
} from "lucide-react";

interface RepoPageProps {
  params: { id: string };
}

function getRelativeTime(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  if (days > 0) return `${days}d ago`;
  if (hours > 0) return `${hours}h ago`;
  if (minutes > 0) return `${minutes}m ago`;
  return "just now";
}

function getEventIcon(eventType: string) {
  switch (eventType) {
    case "push":
      return <GitCommit className="h-3.5 w-3.5 text-green-600" />;
    case "pull_request":
      return <GitPullRequest className="h-3.5 w-3.5 text-purple-600" />;
    case "release":
      return <Tag className="h-3.5 w-3.5 text-blue-600" />;
    default:
      return <GitBranch className="h-3.5 w-3.5 text-gray-600" />;
  }
}

export default function RepoDetailPage({ params }: RepoPageProps) {
  // Fetch showcase repo details from the user's showcase list
  const { data: repos, isLoading } = useQuery<ShowcaseRepo[]>({
    queryKey: ["showcaseRepos"],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: ShowcaseRepo[] }>(
        "/showcase"
      );
      return data.data;
    },
  });

  // Fetch recent activity for this repo from the feed
  const { data: feedData } = useQuery<{ items: ActivityItem[] }>({
    queryKey: ["repoActivity", params.id],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: { items: ActivityItem[] } }>(
        "/feed",
        { params: { limit: 10 } }
      );
      return data.data;
    },
  });

  const repo = repos?.find((r) => r.id === params.id);
  // Filter activities related to this repo
  const repoActivities = (feedData?.items ?? []).filter(
    (item) => repo && item.repo_name.includes(repo.repo_name)
  ).slice(0, 8);

  if (isLoading) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-8 space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-48 w-full rounded-2xl" />
        <Skeleton className="h-64 w-full rounded-2xl" />
      </div>
    );
  }

  if (!repo) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-8">
        <Card>
          <CardContent className="p-12 text-center">
            <div className="h-14 w-14 rounded-full bg-gray-50 flex items-center justify-center mx-auto mb-4">
              <FolderGit2 className="h-7 w-7 text-gray-300" />
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-1">
              Repository not found
            </h3>
            <p className="text-sm text-gray-500 mb-4">
              This repository is not in your showcase or does not exist.
            </p>
            <Link href="/showcase">
              <Button variant="outline" size="sm">
                Back to Showcase
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  const githubUrl = repo.html_url || `https://github.com/${repo.repo_full_name}`;

  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      {/* Back link */}
      <Link
        href="/showcase"
        className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-primary-600 mb-6 transition-colors"
      >
        <ArrowLeft className="h-3.5 w-3.5" />
        Back to Showcase
      </Link>

      {/* Repo Header */}
      <Card className="overflow-hidden mb-6">
        <div className="h-2 gradient-primary" />
        <CardContent className="p-6">
          <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2">
                <FolderGit2 className="h-5 w-5 text-primary-600" />
                <h1 className="text-2xl font-bold text-gray-900">
                  {repo.repo_name}
                </h1>
              </div>
              <p className="text-gray-600 mb-4">
                {repo.description || "No description provided. You can add one from the Showcase page."}
              </p>
              <div className="flex flex-wrap items-center gap-2">
                {repo.language && (
                  <Badge variant="secondary" className="flex items-center gap-1">
                    <Code2 className="h-3 w-3" />
                    {repo.language}
                  </Badge>
                )}
                <Badge className="bg-primary-50 text-primary-700 border-primary-200">
                  {repo.academic_tag.replace("_", " ")}
                </Badge>
                <span className="text-xs text-gray-400">
                  {repo.repo_full_name}
                </span>
              </div>
            </div>
            <div className="flex flex-col gap-2 sm:items-end">
              <a href={githubUrl} target="_blank" rel="noopener noreferrer">
                <Button size="sm" className="gap-1.5 gradient-primary text-white border-0">
                  <ExternalLink className="h-3.5 w-3.5" />
                  View on GitHub
                </Button>
              </a>
              <Link href={`/repos/${repo.id}/discussions`}>
                <Button size="sm" variant="outline" className="gap-1.5">
                  <MessageSquare className="h-3.5 w-3.5" />
                  Discussions
                </Button>
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Content grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Recent Activity (2/3) */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base flex items-center gap-2">
                <GitBranch className="h-4 w-4 text-green-600" />
                Recent Activity
              </CardTitle>
            </CardHeader>
            <CardContent>
              {repoActivities.length > 0 ? (
                <div className="space-y-3">
                  {repoActivities.map((activity) => (
                    <div
                      key={activity.id}
                      className="flex items-start gap-3 p-3 rounded-lg bg-gray-50/50 border border-gray-100"
                    >
                      <div className="mt-0.5">{getEventIcon(activity.event_type)}</div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm text-gray-700 line-clamp-2">
                          {activity.summary}
                        </p>
                        <span className="text-xs text-gray-400 flex items-center gap-1 mt-1">
                          <Clock className="h-3 w-3" />
                          {getRelativeTime(activity.created_at)}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <div className="h-12 w-12 rounded-full bg-gray-50 flex items-center justify-center mx-auto mb-3">
                    <GitBranch className="h-6 w-6 text-gray-300" />
                  </div>
                  <p className="text-sm text-gray-500">
                    No recent activity recorded for this repository.
                  </p>
                  <p className="text-xs text-gray-400 mt-1">
                    Try syncing your GitHub activity from the Feed page.
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Sidebar (1/3) */}
        <div className="space-y-4">
          {/* Repo Info */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm">Repository Info</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-xs text-gray-500">Full Name</span>
                <span className="text-xs font-medium text-gray-700">{repo.repo_full_name}</span>
              </div>
              {repo.language && (
                <div className="flex items-center justify-between">
                  <span className="text-xs text-gray-500">Language</span>
                  <span className="text-xs font-medium text-gray-700">{repo.language}</span>
                </div>
              )}
              <div className="flex items-center justify-between">
                <span className="text-xs text-gray-500">Tag</span>
                <Badge className="text-[10px] bg-primary-50 text-primary-700">
                  {repo.academic_tag.replace("_", " ")}
                </Badge>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-xs text-gray-500">Added</span>
                <span className="text-xs text-gray-700">
                  {new Date(repo.created_at).toLocaleDateString()}
                </span>
              </div>
            </CardContent>
          </Card>

          {/* Quick Actions */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm">Quick Actions</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <a
                href={githubUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 text-sm text-gray-600 hover:text-primary-600 transition-colors"
              >
                <ExternalLink className="h-4 w-4" />
                Open on GitHub
              </a>
              <a
                href={`${githubUrl}/issues`}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 text-sm text-gray-600 hover:text-primary-600 transition-colors"
              >
                <MessageSquare className="h-4 w-4" />
                GitHub Issues
              </a>
              <a
                href={`${githubUrl}/pulls`}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 text-sm text-gray-600 hover:text-primary-600 transition-colors"
              >
                <GitPullRequest className="h-4 w-4" />
                Pull Requests
              </a>
              <Link
                href={`/repos/${repo.id}/discussions`}
                className="flex items-center gap-2 text-sm text-gray-600 hover:text-primary-600 transition-colors"
              >
                <MessageSquare className="h-4 w-4" />
                Forum Discussions
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
