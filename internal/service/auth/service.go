package auth

import (
	"casino_backend/internal/repository"
)

type serv struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
}

func NewService() *serv {
	return &serv{}
}
