package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CommunityHandler handles community-wide HTTP requests.
type CommunityHandler struct {
	pool *pgxpool.Pool
}

// NewCommunityHandler creates a new community handler.
func NewCommunityHandler(pool *pgxpool.Pool) *CommunityHandler {
	return &CommunityHandler{pool: pool}
}

// CommunityStats holds platform-wide statistics.
type CommunityStats struct {
	TotalMembers    int      `json:"total_members"`
	TotalRepos      int      `json:"total_repos"`
	TotalActivities int      `json:"total_activities"`
	ActiveToday     int      `json:"active_today"`
	TopLanguages    []string `json:"top_languages"`
	CommitsThisWeek int      `json:"commits_this_week"`
}

// HandleGetStats handles GET /api/stats.
func (h *CommunityHandler) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats, err := h.getCommunityStats(ctx)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to fetch stats")
		return
	}
	RespondJSON(w, http.StatusOK, stats)
}

func (h *CommunityHandler) getCommunityStats(ctx context.Context) (*CommunityStats, error) {
	stats := &CommunityStats{}

	// Total members
	err := h.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND is_active = true`).Scan(&stats.TotalMembers)
	if err != nil {
		return nil, err
	}

	// Total showcase repos
	err = h.pool.QueryRow(ctx, `SELECT COUNT(*) FROM showcase_repos WHERE deleted_at IS NULL`).Scan(&stats.TotalRepos)
	if err != nil {
		return nil, err
	}

	// Total activities
	err = h.pool.QueryRow(ctx, `SELECT COUNT(*) FROM activity_logs`).Scan(&stats.TotalActivities)
	if err != nil {
		return nil, err
	}

	// Active today (users with activity today)
	today := time.Now().Truncate(24 * time.Hour)
	err = h.pool.QueryRow(ctx, `SELECT COUNT(DISTINCT user_id) FROM activity_logs WHERE created_at >= $1`, today).Scan(&stats.ActiveToday)
	if err != nil {
		return nil, err
	}

	// Commits this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	err = h.pool.QueryRow(ctx, `SELECT COUNT(*) FROM activity_logs WHERE event_type = 'push' AND created_at >= $1`, weekAgo).Scan(&stats.CommitsThisWeek)
	if err != nil {
		return nil, err
	}

	// Top languages from showcase repos
	rows, err := h.pool.Query(ctx, `
		SELECT language, COUNT(*) as cnt 
		FROM showcase_repos 
		WHERE deleted_at IS NULL AND language != '' 
		GROUP BY language 
		ORDER BY cnt DESC 
		LIMIT 8`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var lang string
		var cnt int
		if err := rows.Scan(&lang, &cnt); err != nil {
			continue
		}
		stats.TopLanguages = append(stats.TopLanguages, lang)
	}
	if stats.TopLanguages == nil {
		stats.TopLanguages = []string{}
	}

	return stats, nil
}

// PopularRepo holds a popular showcase repo with activity count.
type PopularRepo struct {
	ID            string `json:"id"`
	RepoName      string `json:"repo_name"`
	RepoFullName  string `json:"repo_full_name"`
	Description   string `json:"description"`
	Language      string `json:"language"`
	HTMLURL       string `json:"html_url"`
	AcademicTag   string `json:"academic_tag"`
	OwnerAlias    string `json:"owner_alias"`
	OwnerAvatar   string `json:"owner_avatar"`
	ActivityCount int    `json:"activity_count"`
}

// HandleGetPopularRepos handles GET /api/repos/popular.
func (h *CommunityHandler) HandleGetPopularRepos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := `
		SELECT s.id, s.repo_name, s.repo_full_name, s.description, s.language, 
			s.html_url, s.academic_tag, u.alias, u.avatar_url,
			COALESCE((SELECT COUNT(*) FROM activity_logs a WHERE a.showcase_repo_id = s.id), 0) as activity_count
		FROM showcase_repos s
		JOIN users u ON s.user_id = u.id
		WHERE s.deleted_at IS NULL
		ORDER BY activity_count DESC, s.created_at DESC
		LIMIT 6`

	rows, err := h.pool.Query(ctx, query)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "failed to fetch popular repos")
		return
	}
	defer rows.Close()

	var repos []PopularRepo
	for rows.Next() {
		var repo PopularRepo
		if err := rows.Scan(
			&repo.ID, &repo.RepoName, &repo.RepoFullName, &repo.Description,
			&repo.Language, &repo.HTMLURL, &repo.AcademicTag, &repo.OwnerAlias,
			&repo.OwnerAvatar, &repo.ActivityCount,
		); err != nil {
			continue
		}
		if repo.HTMLURL == "" {
			repo.HTMLURL = "https://github.com/" + repo.RepoFullName
		}
		repos = append(repos, repo)
	}

	if repos == nil {
		repos = []PopularRepo{}
	}

	RespondJSON(w, http.StatusOK, repos)
}
