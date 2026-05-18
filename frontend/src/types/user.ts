export type Role = "student" | "faculty" | "admin";

export interface User {
  id: string;
  nim: string;
  alias: string;
  bio: string;
  avatar_url: string;
  github_username: string;
  role: Role;
  created_at: string;
  updated_at: string;
}

export interface PublicProfile {
  id: string;
  alias: string;
  bio: string;
  avatar_url: string;
  github_username: string;
  role: Role;
  showcase_repos: import("./showcase").ShowcaseRepo[];
  created_at: string;
}

export interface AcademicIdentity {
  nim: string;
  full_name: string;
  major: string;
  semester: number;
}
