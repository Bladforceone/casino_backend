package line

import (
	"casino_backend/internal/config"
	"casino_backend/internal/repository"
	"casino_backend/internal/service"
)

type serv struct {
	cfg  config.LineConfig
	repo repository.LineRepository
}

// NewLine Создать новый слот 5x3
func NewLineService(cfg config.LineConfig, repo repository.LineRepository) service.LineService {
	return &serv{
		cfg:  cfg,
		repo: repo,
	}
}
