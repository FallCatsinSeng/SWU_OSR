import type { ActivityItem } from "@/types/activity";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { GitBranch, GitPullRequest, Tag } from "lucide-react";

interface ActivityCardProps {
  item: ActivityItem;
}

function getEventIcon(eventType: string) {
  switch (eventType) {
    case "push":
      return <GitBranch className="h-4 w-4 text-green-600" />;
    case "pull_request":
      return <GitPullRequest className="h-4 w-4 text-purple-600" />;
    case "release":
      return <Tag className="h-4 w-4 text-blue-600" />;
    default:
      return <GitBranch className="h-4 w-4 text-gray-600" />;
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
  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          <Avatar
            src={item.avatar_url}
            alt={item.user_alias}
            fallback={item.user_alias.charAt(0).toUpperCase()}
            size="sm"
          />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              {getEventIcon(item.event_type)}
              <span className="font-medium text-sm text-gray-900">
                {item.user_alias}
              </span>
              <Badge variant="secondary" className="text-xs">
                {item.repo_name}
              </Badge>
              <span className="text-xs text-gray-400 ml-auto">
                {getRelativeTime(item.created_at)}
              </span>
            </div>
            <p className="text-sm text-gray-600 truncate">{item.summary}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
