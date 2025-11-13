package line

import (
	"casino_backend/internal/model"
)

type Slot5x3 struct {
	cfg            *model.SlotsConfig
	Bet            int
	Balance        int
	FreeSpins      int
	InFreeSpins    bool
	GuaranteeBonus bool
	LastBoard      [5][3]rune
	reelPool       []string
	reelWeights    []int
}

// NewSlot5x3 Создать новый слот 5x3
func NewSlot5x3(cfg *model.SlotsConfig, balance int) *Slot5x3 {
	syms, weights := buildReelWeights(cfg)
	return &Slot5x3{
		cfg:         cfg,
		Bet:         1,
		Balance:     balance,
		reelPool:    syms, // оставим для совместимости
		reelWeights: weights,
	}
}
