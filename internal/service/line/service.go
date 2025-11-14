package line

import (
	"casino_backend/internal/config"
	"casino_backend/internal/repository"
)

type LineService struct {
	cfg      config.LineConfig
	userRepo repository.UserRepository
	lineRepo repository.LineRepository
}

// NewLine Создать новый слот 5x3
func NewLineService(cfg config.LineConfig, userRepo repository.UserRepository, lineRepo repository.LineRepository) *LineService {
	return &LineService{
		cfg:      cfg,
		userRepo: userRepo,
		lineRepo: lineRepo,
	}
}
