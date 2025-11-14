package line

import (
	"casino_backend/internal/config"
)

type LineService struct {
	cfg *config.LineConfig
}

// NewLine Создать новый слот 5x3
func NewLineService(cfg *config.LineConfig) *LineService {
	return &LineService{
		cfg: cfg,
	}
}
