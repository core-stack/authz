package session

import "context"

func (s *SessionManager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	return s.store.Get(ctx, sessionID)
}

func (s *SessionManager) GetSessionsByUserID(ctx context.Context, UserID string) ([]*Session, error) {
	return s.store.GetByFilter(ctx, map[string]any{"user_id": UserID}, 0, 0)
}

func (s *SessionManager) ListSessions(ctx context.Context, limit, offset int) ([]*Session, error) {
	return s.store.GetByFilter(ctx, map[string]any{}, limit, offset)
}
