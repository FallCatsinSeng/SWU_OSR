export interface LeaderboardEntry {
  rank: number;
  user_id: string;
  alias: string;
  avatar_url: string;
  total_points: number;
  push_points: number;
  pr_points: number;
  forum_points: number;
  other_points: number;
  streak_days: number;
}

export interface LeaderboardResult {
  period: LeaderboardPeriod;
  from: string;
  to: string;
  entries: LeaderboardEntry[];
}

export interface UserPointsSummary {
  user_id: string;
  total_points: number;
  push_count: number;
  pr_count: number;
  thread_count: number;
  comment_count: number;
  showcase_count: number;
  streak_days: number;
  rank: number;
}

export type LeaderboardPeriod = "weekly" | "monthly" | "semester" | "all_time";
