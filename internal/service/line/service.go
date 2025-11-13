package line

import (
	"casino_backend/internal/model"
)

type LineService struct {
	cfg            *model.LineConfig
	Bet            int
	Balance        int
	FreeSpins      int
	InFreeSpins    bool
	GuaranteeBonus bool
	LastBoard      [5][3]rune
	reelPool       []string
	reelWeights    []int
}

// NewLine Создать новый слот 5x3
func NewLineService(cfg *model.LineConfig, balance int) *LineService {
	syms, weights := buildReelWeights(cfg)
	return &LineService{
		cfg:         cfg,
		Bet:         1,
		Balance:     balance,
		reelPool:    syms, // оставим для совместимости
		reelWeights: weights,
	}
}
