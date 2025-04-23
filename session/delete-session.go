package session

import "context"

func (s *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	return s.store.Delete(ctx, sessionID)
}
