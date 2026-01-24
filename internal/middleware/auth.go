package middleware

import (
	"casino_backend/internal/repository"
	"casino_backend/pkg/token"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	CtxUserIDKey    contextKey = "user_id"
	CtxSessionIDKey contextKey = "session_id"
)

type AuthMiddleware struct {
	jwtSecret []byte
	authRepo  repository.AuthRepository // для проверки session_id (опционально)
}

func NewAuthMiddleware(jwtSecret []byte, authRepo repository.AuthRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		authRepo:  authRepo,
	}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, prefix)

		// 2. Verify access token
		claims, err := token.VerifyAccessToken(tokenStr, m.jwtSecret)
		if err != nil {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		// 3. Check session (critical!)
		_, err = m.authRepo.GetUserIDBySessionID(
			r.Context(),
			claims.SessionID,
		)
		if err != nil {
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}

		// 4. Put data into context
		ctx := context.WithValue(r.Context(), CtxUserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, CtxSessionIDKey, claims.SessionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(CtxUserIDKey).(int)
	return id, ok
}
