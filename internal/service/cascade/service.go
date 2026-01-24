package cascade

import (
	"casino_backend/internal/config"
	"casino_backend/internal/repository"
	"casino_backend/internal/service"
)

type serv struct {
	cfg              []config.CascadeConfig
	cascadeRepo      repository.CascadeRepository
	userRepo         repository.UserRepository
	cascadeStatsRepo repository.CascadeStatsRepository
}

// NewCascadeService Создать новый cascade
func NewCascadeService(
	cfg []config.CascadeConfig,
	repo repository.CascadeRepository,
	userRepo repository.UserRepository,
	cascadeStatsRepo repository.CascadeStatsRepository,
) service.CascadeService {
	return &serv{
		cfg:              cfg,
		cascadeRepo:      repo,
		userRepo:         userRepo,
		cascadeStatsRepo: cascadeStatsRepo,
	}
}
