package zmodel

import (
	"database/sql"
	"time"
)

type UserStatus int

const (
	Active UserStatus = iota
	Pending
	Disabled
)

type User struct {
	ID string

	Username string
	Name     string
	Password string
	Email    string

	RoleID int

	Status  UserStatus
	IsAdmin bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CodeType int

const (
	ResetPassword CodeType = iota
	ActiveAccount
	UpdateEmail
)

type Code struct {
	ID        string
	UserID    string
	Token     string
	Type      CodeType
	ExpiresAt time.Time
	CreatedAt time.Time
	UsedAt    sql.NullTime
}

type OAuth2Provider string

const (
	Google OAuth2Provider = "google"
)

type OAuth2User struct {
	ID         string
	UserID     string
	Provider   string
	ProviderID string
	LinkedAt   time.Time

	User *User
}
