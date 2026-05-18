export interface LoginInput {
  nim: string;
  password: string;
}

export interface PendingSession {
  session_id: string;
  github_oauth_url: string;
}

export interface AuthResult {
  access_token: string;
  refresh_token: string;
  user: import("./user").User;
}
