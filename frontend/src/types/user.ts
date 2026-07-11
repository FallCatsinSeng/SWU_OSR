export type Role = 'student' | 'faculty';

export interface User {
  id: string;
  nim: string;
  alias: string;
  bio: string;
  avatar_url: string;
  banner_url: string;
  github_username: string;
  role: Role;
  created_at: string;
  updated_at: string;
}

export interface UserStats {
  total_commits: number;
  total_repos: number;
  languages: string[];
  active_days: number;
  current_streak: number;
  contribution_days: Record<string, number>;
  // Behavioral fields (used for badge computation)
  night_owl_count?: number; // pushes between 00:00–04:00
  early_bird_count?: number; // pushes between 04:00–07:00
  weekend_count?: number; // pushes on Sat/Sun
  total_push_count?: number; // all pushes (for ratio)
  forum_total?: number; // threads + comments
  total_merged_prs?: number; // merged pull requests
  total_showcase_repos?: number; // showcase repos count
}

export interface PublicProfile {
  id: string;
  alias: string;
  bio: string;
  avatar_url: string;
  banner_url: string;
  github_username: string;
  role: Role;
  showcase_repos: import('./showcase').ShowcaseRepo[];
  stats: UserStats | null;
  created_at: string;
}

export interface AcademicIdentity {
  nim: string;
  full_name: string;
  major: string;
  semester: number;
}
