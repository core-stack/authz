package session

import (
	"context"
	"time"

	"github.com/core-stack/authz/jwt"
	"github.com/core-stack/authz/store"
	"github.com/core-stack/authz/zmodel"
)

var DefaultSessionDuration = 15 * time.Minute

type Session struct {
	ID        string
	UserID    string
	Email     string
	Name      string
	RoleID    int
	Extra     map[string]string
	ExpiresAt time.Time
}

func (s Session) GetID() string {
	return s.ID
}
func (s Session) GetExpiresAt() time.Time {
	return s.ExpiresAt
}

type ISessionManager interface {
	CreateSession(ctx context.Context, user zmodel.User, extra map[string]string) (*jwt.TokenPair, error)
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	GetSessionsByUserID(ctx context.Context, UserID string) ([]*Session, error)
	ListSessions(ctx context.Context, limit, offset int) ([]*Session, error)
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, sessionID string) error
}

type SessionManager struct {
	store store.Store[Session]

	sessionDuration time.Duration
	jwtService      *jwt.JWTService
}

func NewSessionManager(store store.Store[Session], sessionDuration time.Duration, jwtService *jwt.JWTService) *SessionManager {
	return &SessionManager{
		store:           store,
		sessionDuration: sessionDuration,
		jwtService:      jwtService,
	}
}
