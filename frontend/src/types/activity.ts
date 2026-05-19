export type EventType = "push" | "pull_request" | "release";

export interface ActivityItem {
  id: string;
  user_id: string;
  user_alias: string;
  avatar_url: string;
  event_type: EventType;
  repo_id: string;
  repo_name: string;
  summary: string;
  metadata: Record<string, unknown>;
  created_at: string;
}

export interface FeedResponse {
  items: ActivityItem[];
  next_cursor: string;
  has_more: boolean;
}
