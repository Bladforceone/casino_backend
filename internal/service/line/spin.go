package line

import (
	"casino_backend/internal/model"
	"errors"
	"math/rand"
)

func (s *Slot5x3) Spin() (model.SpinResult, error) {
	var res model.SpinResult

	// платный или фриспин?
	if s.FreeSpins == 0 {
		if s.Balance < s.Bet {
			return res, errors.New("not enough balance")
		}
		s.Balance -= s.Bet
	} else {
		res.InFreeSpin = true
		s.FreeSpins--
	}

	// делаем спин
	spin := s.spinOnce(false)

	// обновляем баланс
	s.Balance += spin.TotalPayout

	return model.SpinResult{
		Board:            res.Board,
		LineWins:         res.LineWins,
		ScatterCount:     res.ScatterCount,
		ScatterPayout:    res.ScatterPayout,
		AwardedFreeSpins: res.AwardedFreeSpins,
		TotalPayout:      res.TotalPayout,
		Balance:          res.Balance,
		FreeSpinCount:    s.FreeSpins,
		InFreeSpin:       res.InFreeSpin,
	}, nil
}

// ---------- ВСПОМОГАТЕЛЬНОЕ ----------
// Рандомный выбор по весам
func randomWeighted(symbols []string, weights []int) string {
	total := 0
	for _, w := range weights {
		total += w
	}
	if total <= 0 {
		return symbols[0]
	}
	r := rand.Intn(total)
	for i, w := range weights {
		if r < w {
			return symbols[i]
		}
		r -= w
	}
	return symbols[len(symbols)-1]
}

// Применяет лимит по максимальному выигрышу
func applyMaxPayout(amount, bet, maxMult int) int {
	maxPay := maxMult * bet
	if amount > maxPay {
		return maxPay
	}
	return amount
}

// Строим массив символов и их весов для randomWeighted
func buildReelWeights(cfg *model.SlotsConfig) ([]string, []int) {
	var syms []string
	var weights []int
	for sym, w := range cfg.SymbolWeights {
		if sym == "W" {
			continue
		}
		syms = append(syms, sym)
		weights = append(weights, w)
	}
	return syms, weights
}

// Генерация игрового поля матрицы 5x3
func (s *Slot5x3) generateBoard(forceBonus bool) [5][3]rune {
	var board [5][3]rune

	wildReels := map[int]bool{}
	for r := 1; r <= 3; r++ {
		if rand.Float64() < s.cfg.WildChanceOnReel234 || s.InFreeSpins {
			wildReels[r] = true
		}
	}

	scatterPositions := map[[2]int]bool{}
	if forceBonus || s.GuaranteeBonus {
		s.GuaranteeBonus = false
		cols := rand.Perm(s.cfg.Reels)[:3]
		for _, reel := range cols {
			row := rand.Intn(s.cfg.Rows)
			scatterPositions[[2]int{reel, row}] = true
		}
	}

	for r := 0; r < s.cfg.Reels; r++ {
		for row := 0; row < s.cfg.Rows; row++ {
			if wildReels[r] {
				board[r][row] = 'W'
			} else if scatterPositions[[2]int{r, row}] {
				board[r][row] = 'B'
			} else {
				sym := randomWeighted(s.reelPool, s.reelWeights)
				board[r][row] = rune(sym[0])
			}
		}
	}

	return board
}

// ---------- ОЦЕНКА ЛИНИЙ ----------
func (s *Slot5x3) evaluateLines(board [5][3]rune) []model.LineWin {
	var wins []model.LineWin
	for i, line := range s.cfg.Paylines {
		symbols := make([]rune, s.cfg.Reels)
		for r := 0; r < s.cfg.Reels; r++ {
			symbols[r] = board[r][line[r]]
		}
		if symbols[0] == 'B' {
			continue
		}
		var base rune
		for _, sym := range symbols {
			if sym != 'W' && sym != 'B' {
				base = sym
				break
			}
		}

		if base == 0 {
			continue
		}
		count := 0
		for _, sym := range symbols {
			if sym == base || sym == 'W' {
				count++
			} else {
				break
			}
		}
		if count >= 3 {
			if pays, ok := s.cfg.Pays[string(base)]; ok {
				if val, ok := pays[count]; ok {
					win := model.LineWin{
						Line:   i + 1,
						Symbol: base,
						Count:  count,
						Payout: val * s.Bet / 100,
					}
					wins = append(wins, win)
				}
			}
		}
	}
	return wins
}

// ---------- СПИН (возвращает единый SpinResult) ----------
func (s *Slot5x3) spinOnce(forceBonus bool) model.SpinResult {
	board := s.generateBoard(forceBonus)
	s.LastBoard = board

	// count scatters
	scatters := 0
	for r := 0; r < s.cfg.Reels; r++ {
		for c := 0; c < s.cfg.Rows; c++ {
			if board[r][c] == 'B' {
				scatters++
			}
		}
	}

	// scatter payout
	var scatterPayout int
	if scatters > 0 {
		if val, ok := s.cfg.Pays["B"][scatters]; ok {
			scatterPayout = val * s.Bet / 100
		}
	}

	// line wins
	lineWins := s.evaluateLines(board)
	var lineTotal int
	for _, w := range lineWins {
		lineTotal += w.Payout
	}

	total := applyMaxPayout(lineTotal+scatterPayout, s.Bet, s.cfg.MaxPayoutMultiplier)

	awarded := 0
	if scatters >= 3 {
		if v, ok := s.cfg.FreeSpinsByScatter[scatters]; ok {
			awarded = v
			s.FreeSpins += v
		}
	}

	return model.SpinResult{
		Board:            board,
		LineWins:         lineWins,
		ScatterCount:     scatters,
		ScatterPayout:    scatterPayout,
		AwardedFreeSpins: awarded,
		TotalPayout:      total,
	}
}
