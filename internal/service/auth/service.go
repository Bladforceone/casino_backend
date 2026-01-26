package auth

import (
	"casino_backend/internal/config"
	"casino_backend/internal/repository"
	"casino_backend/internal/service"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
)

// Проверка соответствия интерфейсу
var _ service.AuthService = (*serv)(nil)

type serv struct {
	txManager trm.Manager
	jwtConfig config.JWTConfig
	userRepo  repository.UserRepository
	authRepo  repository.AuthRepository
}

func NewService(
	txManager trm.Manager,
	jwtConfig config.JWTConfig,
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
) *serv {
	return &serv{
		txManager: txManager,
		jwtConfig: jwtConfig,
		userRepo:  userRepo,
		authRepo:  authRepo,
	}
}

func generateSessionID() string {
	return uuid.New().String()
}
