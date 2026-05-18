package domain

import (
	"time"

	"github.com/google/uuid"
)

// AcademicTag categorizes a showcase repository.
type AcademicTag string

const (
	TagCoursework       AcademicTag = "coursework"
	TagThesis           AcademicTag = "thesis"
	TagHackathon        AcademicTag = "hackathon"
	TagPersonalResearch AcademicTag = "personal_research"
	TagTeamProject      AcademicTag = "team_project"
)

// ShowcaseRepo represents a student's selected repository for public showcase.
type ShowcaseRepo struct {
	ID           uuid.UUID   `json:"id"`
	UserID       uuid.UUID   `json:"user_id"`
	GitHubRepoID int64       `json:"github_repo_id"`
	RepoName     string      `json:"repo_name"`
	RepoFullName string      `json:"repo_full_name"`
	Description  string      `json:"description"`
	Language     string      `json:"language"`
	HTMLURL      string      `json:"html_url"`
	AcademicTag  AcademicTag `json:"academic_tag"`
	WebhookID    *int64      `json:"-"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	DeletedAt    *time.Time  `json:"-"`
}

// ShowcaseSelection is the input for selecting repos to showcase.
type ShowcaseSelection struct {
	RepoID   int64       `json:"repo_id"`
	RepoName string      `json:"repo_name"`
	FullName string      `json:"full_name"`
	Tag      AcademicTag `json:"tag"`
}
