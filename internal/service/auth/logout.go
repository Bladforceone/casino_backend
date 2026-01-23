package auth

import "context"

func (s *serv) Logout(ctx context.Context, sessionID string) error {
	return s.authRepo.DeleteSession(ctx, sessionID)
}
