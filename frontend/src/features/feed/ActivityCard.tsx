import Link from "next/link";
import type { ActivityItem } from "@/types/activity";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  GitBranch,
  GitPullRequest,
  Tag,
  GitCommit,
  GitFork,
  Star,
  CircleDot,
  Plus,
} from "lucide-react";

interface ActivityCardProps {
  item: ActivityItem;
}

function getEventConfig(eventType: string, summary: string) {
  if (summary.startsWith("Forked")) {
    return { icon: GitFork, label: "forked" };
  }
  if (summary.startsWith("Starred")) {
    return { icon: Star, label: "starred" };
  }
  if (summary.startsWith("Created repository")) {
    return { icon: Plus, label: "created" };
  }
  if (
    summary.startsWith("Created branch") ||
    summary.startsWith("Created tag")
  ) {
    return { icon: GitBranch, label: "created" };
  }
  if (summary.startsWith("Issue")) {
    return { icon: CircleDot, label: "issue on" };
  }

  switch (eventType) {
    case "push":
      return { icon: GitCommit, label: "pushed to" };
    case "pull_request":
      return { icon: GitPullRequest, label: "PR on" };
    case "release":
      return { icon: Tag, label: "released" };
    default:
      return { icon: GitBranch, label: "updated" };
  }
}

function getRelativeTime(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  const weeks = Math.floor(days / 7);

  if (weeks > 0) return `${weeks}w ago`;
  if (days > 0) return `${days}d ago`;
  if (hours > 0) return `${hours}h ago`;
  if (minutes > 0) return `${minutes}m ago`;
  return "just now";
}

function getRepoShortName(repoName: string): string {
  if (repoName.includes("/")) {
    return repoName.split("/").pop() || repoName;
  }
  return repoName;
}

export function ActivityCard({ item }: ActivityCardProps) {
  const config = getEventConfig(item.event_type, item.summary);
  const Icon = config.icon;
  const repoShort = getRepoShortName(item.repo_name);

  return (
    <Card className="transition-shadow hover:shadow-geist-3">
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          {/* Event icon — neutral canvas-soft-2 background */}
          <div className="h-8 w-8 rounded-geist-sm bg-geist-canvas-soft-2 dark:bg-neutral-800 flex items-center justify-center shrink-0">
            <Icon className="h-4 w-4 text-geist-ink dark:text-neutral-50" />
          </div>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <Link
                href={`/profiles/${item.user_alias}`}
                className="text-body-sm-strong text-geist-ink dark:text-neutral-50 hover:text-geist-link dark:hover:text-cyan-400 transition-colors"
              >
                {item.user_alias}
              </Link>
              <span className="text-caption text-geist-mute dark:text-neutral-500">
                {config.label}
              </span>
              <Badge variant="secondary" className="text-[10px]">
                {repoShort}
              </Badge>
            </div>
            <p className="text-body-sm text-geist-body dark:text-neutral-400 mt-1 line-clamp-2">
              {item.summary}
            </p>
          </div>

          {/* Timestamp — mono caption */}
          <span className="text-caption-mono text-geist-mute dark:text-neutral-500 shrink-0 pt-0.5">
            {getRelativeTime(item.created_at)}
          </span>
        </div>
      </CardContent>
    </Card>
  );
}
