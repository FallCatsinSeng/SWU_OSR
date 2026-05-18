package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user's role in the system.
type Role string

const (
	RoleStudent Role = "student"
	RoleFaculty Role = "faculty"
	RoleAdmin   Role = "admin"
)

// User is the central identity binding entity.
type User struct {
	ID             uuid.UUID  `json:"id"`
	NIM            string     `json:"nim"`
	FullName       string     `json:"-"`
	Major          string     `json:"-"`
	Semester       int        `json:"-"`
	Alias          string     `json:"alias"`
	Bio            string     `json:"bio"`
	AvatarURL      string     `json:"avatar_url"`
	GitHubUsername string     `json:"github_username"`
	GitHubID       int64      `json:"-"`
	GitHubToken    string     `json:"-"`
	Role           Role       `json:"role"`
	IsActive       bool       `json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"-"`
}

// UserClaims represents JWT claims stored in request context.
type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Alias  string    `json:"alias"`
	Role   Role      `json:"role"`
}
