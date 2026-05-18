# Requirements Document

## Introduction

SWU OSR (Open Source Repository) is a campus open-source repository platform for SWU (Sains dan Teknologi Walisongo) university in Purwokerto. The platform links student academic identities from SIAKAD with their GitHub accounts, aggregates GitHub activity via webhooks, and displays it on a public dashboard with a pseudonym-first culture. It includes an internal discussion forum for campus-scoped collaboration. The system enforces strict identity privacy where only faculty and admin roles can view real academic identities behind student aliases.

## Glossary

- **Platform**: The SWU OSR web application (backend + frontend)
- **SIAKAD_Proxy**: The backend component that authenticates students by forwarding credentials to the campus SIAKAD server without storing them
- **GitHub_OAuth_Client**: The component handling GitHub OAuth authorization flow and API access
- **Auth_Service**: The service orchestrating two-step authentication (SIAKAD + GitHub) and token management
- **Showcase_Service**: The service managing student repository selection, academic tagging, and webhook registration
- **Aggregator_Service**: The service processing incoming GitHub webhooks and storing activity logs
- **Profile_Service**: The service managing public profiles and role-gated identity access
- **Forum_Service**: The service managing discussion threads, comments, and notifications
- **Student**: A user with role "student" who has completed SIAKAD + GitHub identity binding
- **Faculty**: A user with role "faculty" who can view real academic identities
- **Admin**: A user with role "admin" who can view real identities and perform manual verification
- **NIM**: Student identification number (Nomor Induk Mahasiswa) from SIAKAD
- **Alias**: A pseudonym chosen by the student for public display
- **Showcase_Repo**: A GitHub repository selected by a student for display on the platform
- **Academic_Tag**: A classification label for showcase repositories (coursework, thesis, hackathon, personal_research, team_project)
- **Activity_Log**: A record of a GitHub event (push, pull_request, release) captured via webhook
- **Webhook_Secret**: A shared HMAC-SHA256 secret used to verify GitHub webhook payloads
- **JWT**: JSON Web Token used for stateless authentication (15-minute expiry)
- **Refresh_Token**: A single-use token for obtaining new JWTs without re-authentication

## Requirements

### Requirement 1: SIAKAD Authentication

**User Story:** As a student, I want to log in using my SIAKAD credentials, so that my academic identity is verified without the platform storing my password.

#### Acceptance Criteria

1. WHEN a student submits NIM and password, THE SIAKAD_Proxy SHALL forward the credentials to the campus SIAKAD server and return student data (full name, NIM, major, semester, active status) upon successful validation
2. WHEN the SIAKAD server rejects the credentials, THE SIAKAD_Proxy SHALL return an authentication error without storing any credential data
3. THE Platform SHALL never persist student SIAKAD passwords in any database table, log file, or cache
4. WHEN the SIAKAD server is unavailable (timeout or 5xx response), THE Auth_Service SHALL retry with exponential backoff up to 3 attempts and then return a 503 error to the client
5. IF the SIAKAD server is persistently unavailable, THEN THE Auth_Service SHALL allow admin users to manually verify students via the ManualVerify endpoint

### Requirement 2: GitHub OAuth Identity Binding

**User Story:** As a student, I want to link my GitHub account after SIAKAD login, so that my academic identity is bound to my GitHub username.

#### Acceptance Criteria

1. WHEN SIAKAD authentication succeeds, THE Auth_Service SHALL create a pending session and redirect the student to the GitHub OAuth authorization page
2. WHEN GitHub returns an authorization code, THE GitHub_OAuth_Client SHALL exchange it for an access token and retrieve the GitHub username
3. WHEN the GitHub OAuth flow completes, THE Auth_Service SHALL create a user record binding student_id, github_username, and encrypted access token
4. THE Platform SHALL enforce uniqueness constraints: each NIM maps to exactly one GitHub username, and each GitHub username maps to exactly one NIM
5. WHEN identity binding completes successfully, THE Auth_Service SHALL issue a JWT (15-minute expiry) and a single-use refresh token stored as a hash in the database; tokens SHALL NOT be created if binding fails

### Requirement 3: Token Lifecycle Management

**User Story:** As an authenticated user, I want my session to be managed securely, so that my access is protected and I can stay logged in without re-authenticating frequently.

#### Acceptance Criteria

1. THE Auth_Service SHALL issue JWTs with a 15-minute expiration containing user ID, role, and alias claims
2. WHEN a JWT has expired, THE Auth_Service SHALL allow the client to exchange a valid refresh token for a new JWT and refresh token pair; preemptive refresh before expiration SHALL NOT be supported
3. WHEN a refresh token is used, THE Auth_Service SHALL revoke it immediately so it cannot be reused
4. WHEN a user logs out, THE Auth_Service SHALL revoke all refresh tokens associated with that user
5. WHEN a request contains an invalid or expired JWT, THE Platform SHALL respond with 401 Unauthorized without invoking the handler

### Requirement 4: Profile and Onboarding

**User Story:** As a student, I want to set up a pseudonymous public profile, so that I can participate in the campus open-source community with a chosen identity.

#### Acceptance Criteria

1. WHEN a student completes identity binding, THE Platform SHALL direct them to the onboarding flow to set alias, bio, and avatar
2. THE Profile_Service SHALL validate aliases as 3-50 characters, alphanumeric plus underscores, and unique across all users
3. WHEN a student updates their profile, THE Profile_Service SHALL persist the alias, bio, and avatar URL to the database
4. THE Platform SHALL display only the alias and avatar on public-facing pages, never the real name or NIM

### Requirement 5: Repository Showcase Selection

**User Story:** As a student, I want to select which GitHub repositories to showcase on the platform, so that I can curate my campus-visible portfolio.

#### Acceptance Criteria

1. WHEN a student requests available repositories, THE Showcase_Service SHALL fetch the student's public repositories from GitHub using the stored access token
2. WHEN a student submits showcase selections (max 20 repositories), THE Showcase_Service SHALL store each selection with its academic tag in the database
3. THE Showcase_Service SHALL validate that each academic tag is one of: coursework, thesis, hackathon, personal_research, or team_project
4. WHEN showcase selections are saved, THE Showcase_Service SHALL register GitHub webhooks for push, pull_request, and release events on each selected repository
5. WHEN a repository is removed from the showcase, THE Showcase_Service SHALL remove the corresponding GitHub webhook
6. WHEN showcase selections are updated, THE Showcase_Service SHALL perform the operation atomically: all changes succeed or none are applied

### Requirement 6: Webhook Processing and Activity Logging

**User Story:** As a student, I want my GitHub activity to be automatically captured, so that my contributions are visible on the campus dashboard without manual action.

#### Acceptance Criteria

1. WHEN a GitHub webhook is received, THE Aggregator_Service SHALL verify the HMAC-SHA256 signature before processing the payload
2. IF the webhook signature is invalid, THEN THE Aggregator_Service SHALL reject the payload and perform no database writes
3. WHEN a valid push event is received from a registered user and showcased repository, THE Aggregator_Service SHALL create an activity log entry with commit count, branch reference, and up to 5 commit details
4. WHEN a valid pull_request event is received from a registered user and showcased repository, THE Aggregator_Service SHALL create an activity log entry with action, PR number, and title
5. WHEN a valid release event is received from a registered user and showcased repository, THE Aggregator_Service SHALL create an activity log entry with tag name and release name
6. WHEN a webhook is received from a non-registered user or non-showcased repository, THE Aggregator_Service SHALL silently ignore it and return 200 to GitHub
7. WHEN a duplicate webhook event is received (same github_event_id), THE Aggregator_Service SHALL ignore the duplicate and not create a second activity log entry

### Requirement 7: Activity Feed and Pagination

**User Story:** As a visitor, I want to browse the campus activity feed, so that I can see recent open-source contributions from students.

#### Acceptance Criteria

1. THE Aggregator_Service SHALL return activity feed items ordered by creation time in strictly descending order
2. THE Aggregator_Service SHALL support cursor-based pagination with a configurable limit (1-50 items per page, default 20)
3. THE Aggregator_Service SHALL always include a next_cursor field in paginated responses; WHEN more items exist beyond the current page, next_cursor SHALL contain the cursor value for the subsequent page; WHEN no more pages exist, next_cursor SHALL be an empty string
4. WHEN all pages are iterated, THE Aggregator_Service SHALL have returned every activity log entry exactly once with no duplicates and no gaps (assuming static data)
5. IF cursor generation fails during a paginated request, THEN THE Aggregator_Service SHALL fail the entire request with an error rather than returning partial results without a cursor

### Requirement 8: Public Profile and Statistics

**User Story:** As a visitor, I want to view student profiles with their showcase repositories and contribution statistics, so that I can discover active contributors.

#### Acceptance Criteria

1. WHEN a visitor requests a public profile by alias, THE Profile_Service SHALL return the alias, bio, avatar, showcase repositories, contribution statistics, badges, and join date only for users with at least one commit or showcase repository
2. THE Profile_Service SHALL compute user statistics including total commits, total repositories, programming languages used, active days, and current streak
3. THE Profile_Service SHALL never include NIM, full name, major, or semester in public profile responses

### Requirement 9: Role-Gated Identity Access

**User Story:** As a faculty member, I want to view the real academic identity behind a student alias, so that I can identify students for academic purposes.

#### Acceptance Criteria

1. WHEN a faculty or admin user requests real identity for an alias, THE Profile_Service SHALL return the full name, NIM, major, and semester
2. WHEN a student or unauthenticated user requests real identity, THE Profile_Service SHALL respond with 403 Forbidden
3. THE Platform SHALL only display the "View Real Identity" button to users with faculty or admin role in their session

### Requirement 10: Discussion Forum Threads

**User Story:** As a student, I want to create discussion threads on showcase repositories, so that I can collaborate with peers on campus-specific topics.

#### Acceptance Criteria

1. WHEN a user creates a thread on a showcase repository, THE Forum_Service SHALL store the thread with title (5-255 chars), body (1-10000 chars), author ID, and repository ID
2. WHEN threads are listed for a repository, THE Forum_Service SHALL return them ordered by creation time descending with pagination support
3. WHEN a thread is created, THE Forum_Service SHALL create a notification for the repository owner

### Requirement 11: Discussion Forum Comments

**User Story:** As a student, I want to reply to discussion threads, so that I can participate in technical conversations within the campus platform.

#### Acceptance Criteria

1. WHEN a user posts a comment on a thread, THE Forum_Service SHALL store the comment with body (1-10000 chars), author ID, thread ID, and optional parent comment ID for flat replies; comment storage SHALL succeed independently of the count increment
2. WHEN a comment is successfully stored, THE Forum_Service SHALL increment the thread's comment count; IF the count increment fails, THEN THE Forum_Service SHALL allow eventual consistency and not roll back the stored comment
3. WHEN a comment is successfully stored on a thread, THE Forum_Service SHALL create a notification for the thread author

### Requirement 12: Forum Notifications

**User Story:** As a repository owner or thread author, I want to receive notifications about new discussions and replies, so that I can stay informed about activity on my content.

#### Acceptance Criteria

1. THE Forum_Service SHALL create notifications only for the repository owner (on new threads) or the thread author (on new comments)
2. WHEN a user requests their notifications, THE Forum_Service SHALL return them ordered by creation time descending
3. WHEN a user marks a notification as read, THE Forum_Service SHALL update the notification's read status; IF the backend update fails, THEN THE Platform SHALL allow the UI action to succeed and reconcile the read status eventually

### Requirement 13: Forum Data Isolation

**User Story:** As a platform administrator, I want forum discussions to remain internal to the platform, so that academic privacy is maintained and discussions are not exposed to public GitHub.

#### Acceptance Criteria

1. THE Forum_Service SHALL store all thread and comment data exclusively in the SWU OSR internal database
2. THE Forum_Service SHALL never write discussion data to GitHub Issues, Discussions, or any external service
3. THE Forum_Service SHALL never read from GitHub Issues or Discussions to populate forum content
4. IF the internal database is temporarily unavailable, THEN THE Forum_Service SHALL use temporary external storage as a fallback and reconcile data to the internal database when it recovers

### Requirement 14: Input Validation and Security

**User Story:** As a platform administrator, I want all user inputs to be validated and sanitized, so that the system is protected against injection attacks and malformed data.

#### Acceptance Criteria

1. THE Platform SHALL validate all request inputs at the handler layer before passing them to service logic
2. THE Platform SHALL use parameterized queries (via sqlc) for all database operations to prevent SQL injection
3. THE Platform SHALL encrypt GitHub access tokens at rest using AES-256-GCM with an environment-sourced key, and SHALL fail to start if the encryption key environment variable is missing or invalid
4. THE Platform SHALL enforce rate limiting per IP (100 requests/minute) and per user (300 requests/minute) via Redis sliding window
5. THE Platform SHALL restrict CORS to the specific frontend origin only, with no wildcard in production; non-production environments MAY use wildcard origins for development convenience

### Requirement 15: GitHub Token Invalidation Handling

**User Story:** As a student, I want the platform to handle GitHub token expiration gracefully, so that I can re-authorize without losing my profile data.

#### Acceptance Criteria

1. WHEN the GitHub API returns 401 for a stored access token, THE Platform SHALL mark the token as invalid in the database
2. WHEN a token is marked invalid through the GitHub 401 detection process, THE Platform SHALL prompt the user to re-authorize via GitHub OAuth; re-authorization SHALL only be triggered by explicit token invalidation, not by other error conditions
3. WHEN a token is invalidated, THE Platform SHALL continue processing webhooks for the user's repositories (webhooks use app-level secret, not user token)
4. IF the app-level webhook secret is compromised, THEN THE Platform SHALL stop all webhook processing until the secret is rotated and re-registered
