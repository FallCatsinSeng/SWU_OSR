import Link from "next/link";
import type { ActivityItem } from "@/types/activity";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { GitBranch, GitPullRequest, Tag, GitCommit } from "lucide-react";

interface ActivityCardProps {
  item: ActivityItem;
}

function getEventConfig(eventType: string) {
  switch (eventType) {
    case "push":
      return {
        icon: GitCommit,
        color: "text-green-600",
        bg: "bg-green-50",
        label: "pushed to",
      };
    case "pull_request":
      return {
        icon: GitPullRequest,
        color: "text-purple-600",
        bg: "bg-purple-50",
        label: "opened PR in",
      };
    case "release":
      return {
        icon: Tag,
        color: "text-blue-600",
        bg: "bg-blue-50",
        label: "released",
      };
    default:
      return {
        icon: GitBranch,
        color: "text-gray-600",
        bg: "bg-gray-50",
        label: "updated",
      };
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

  if (days > 0) return `${days}d ago`;
  if (hours > 0) return `${hours}h ago`;
  if (minutes > 0) return `${minutes}m ago`;
  return "just now";
}

export function ActivityCard({ item }: ActivityCardProps) {
  const config = getEventConfig(item.event_type);
  const Icon = config.icon;

  return (
    <Card className="hover:border-gray-200 hover:shadow-sm transition-all duration-200 group">
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          {/* Event icon */}
          <div
            className={`h-9 w-9 rounded-lg ${config.bg} flex items-center justify-center shrink-0`}
          >
            <Icon className={`h-4 w-4 ${config.color}`} />
          </div>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <Link
                href={`/profiles/${item.user_alias}`}
                className="font-medium text-sm text-gray-900 hover:text-primary-600 transition-colors"
              >
                {item.user_alias}
              </Link>
              <span className="text-xs text-gray-400">{config.label}</span>
              <Badge
                variant="secondary"
                className="text-[10px] bg-gray-50 text-gray-600 border-gray-100"
              >
                {item.repo_name}
              </Badge>
            </div>
            <p className="text-sm text-gray-600 mt-1 line-clamp-1">
              {item.summary}
            </p>
          </div>

          {/* Time */}
          <span className="text-[11px] text-gray-400 shrink-0 pt-0.5">
            {getRelativeTime(item.created_at)}
          </span>
        </div>
      </CardContent>
    </Card>
  );
}
