package session

import "context"

func (s *SessionManager) UpdateSession(ctx context.Context, session *Session) error {
	return s.store.Set(ctx, *session)
}
