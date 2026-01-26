package cascade

import (
	"casino_backend/internal/config"
	"casino_backend/internal/repository"
	"casino_backend/internal/service"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type serv struct {
	cfg              config.CascadeConfig
	cascadeRepo      repository.CascadeRepository
	userRepo         repository.UserRepository
	cascadeStatsRepo repository.CascadeStatsRepository
	txManager        trm.Manager
}

// NewCascadeService Создать новый cascade
func NewCascadeService(
	cfg config.CascadeConfig,
	repo repository.CascadeRepository,
	userRepo repository.UserRepository,
	cascadeStatsRepo repository.CascadeStatsRepository,
	txManager trm.Manager,
) service.CascadeService {
	return &serv{
		cfg:              cfg,
		cascadeRepo:      repo,
		userRepo:         userRepo,
		cascadeStatsRepo: cascadeStatsRepo,
		txManager:        txManager,
	}
}
