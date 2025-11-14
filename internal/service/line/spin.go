package line

import (
	"casino_backend/internal/model"
	"context"
	"errors"
	"math/rand"
)

var (
	// Линии выплат
	playLines = [][]int{
		{1, 1, 1, 1, 1},
		{0, 0, 0, 0, 0},
		{2, 2, 2, 2, 2},
		{0, 1, 2, 1, 0},
		{2, 1, 0, 1, 2},
		{0, 0, 1, 0, 0},
		{2, 2, 1, 2, 2},
		{1, 0, 0, 0, 1},
		{1, 2, 2, 2, 1},
		{1, 0, 1, 0, 1},
		{1, 2, 1, 2, 1},
		{0, 1, 0, 1, 0},
		{2, 1, 2, 1, 2},
		{1, 1, 0, 1, 1},
		{1, 1, 2, 1, 1},
		{0, 1, 1, 1, 2},
		{2, 1, 1, 1, 0},
		{0, 0, 1, 2, 2},
		{2, 2, 1, 0, 0},
		{1, 0, 2, 0, 1},
	}
	// Таблица выплат
	pays = map[string]map[int]int{
		"S1": {3: 25, 4: 150, 5: 450},
		"S2": {3: 25, 4: 150, 5: 450},
		"S3": {3: 25, 4: 150, 5: 450},
		"S4": {3: 25, 4: 150, 5: 450},
		"S5": {3: 75, 4: 250, 5: 1000},
		"S6": {3: 125, 4: 500, 5: 2500},
		"S7": {3: 125, 4: 500, 5: 2500},
		"S8": {3: 250, 4: 1250, 5: 12500},
		"B":  {1: 50, 2: 100, 3: 250, 4: 500, 5: 1000},
	}
	// Бесплатные вращения
)

const (
	// Барабаны
	reels = 5
	// Линии
	rows = 3

	// Стоимость покупки бонуса (x ставки)
	buyBonusMultiplier = 100
	// Максимальная выплата в кратности ставки
	maxPayoutMultiplier = 10000
)

func (s *LineService) Spin(ctx context.Context, spinReq model.LineSpin) (*model.SpinResult, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, errors.New("user ID not found in context")
	}

	// Получаем текущее количество фриспинов
	countFreeSpins, err := s.lineRepo.GetCountFreeSpins(userID)
	if err != nil {
		return nil, errors.New("failed to get count free spins")
	}
	// Инициализируем результат (чтобы безопасно устанавливать InFreeSpin)
	var res *model.SpinResult

	// платный или фриспин?
	if countFreeSpins == 0 {
		userBalance, err := s.userRepo.GetBalance(userID)
		if err != nil {
			return nil, errors.New("failed to get user balance")
		}
		if userBalance < spinReq.Bet {
			return res, errors.New("not enough balance")
		}

		// Списание должно быть атомарным! Здесь просто пример — в реале делайте в транзакции.
		userBalance -= spinReq.Bet
		if err := s.userRepo.UpdateBalance(userID, userBalance); err != nil {
			return nil, errors.New("failed to update user balance")
		}
	} else {
		// фриспин — уменьшить счётчик сразу
		res = &model.SpinResult{InFreeSpin: true}
		if err := s.lineRepo.UpdateCountFreeSpins(userID, countFreeSpins-1); err != nil {
			return nil, errors.New("failed to update count free spins")
		}
	}

	// делаем спин
	res, err = s.SpinOnce(ctx, false, spinReq)
	if err != nil {
		return nil, err
	}

	// Если это был фриспин, сохраняем флаг
	if countFreeSpins > 0 {
		res.InFreeSpin = true
	}

	// обновляем баланс
	balance, err := s.userRepo.GetBalance(userID)
	if err != nil {
		return nil, errors.New("failed to get user balance")
	}
	balance += res.TotalPayout

	err = s.userRepo.UpdateBalance(userID, balance)
	if err != nil {
		return nil, errors.New("failed to update user balance")
	}

	// Обновляем индекс свободных спинов в возвращаемом результате (актуально)
	freeCount, err := s.lineRepo.GetCountFreeSpins(userID)
	if err != nil {
		return nil, errors.New("failed to get count free spins")
	}

	return &model.SpinResult{
		Board:            res.Board,
		LineWins:         res.LineWins,
		ScatterCount:     res.ScatterCount,
		ScatterPayout:    res.ScatterPayout,
		AwardedFreeSpins: res.AwardedFreeSpins,
		TotalPayout:      res.TotalPayout,
		Balance:          res.Balance,
		FreeSpinCount:    freeCount,
		InFreeSpin:       res.InFreeSpin,
	}, nil
}

// SpinOnce выполняет один спин (возвращает единый SpinResult)
func (s *LineService) SpinOnce(ctx context.Context, forceBonus bool, spinReq model.LineSpin) (*model.SpinResult, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, errors.New("user ID not found in context")
	}

	board, err := s.GenerateBoard(forceBonus, spinReq, userID)
	if err != nil {
		return nil, err
	}

	// count scatters
	scatters := 0
	for r := 0; r < reels; r++ {
		for c := 0; c < rows; c++ {
			if board[r][c] == "B" {
				scatters++
			}
		}
	}

	// scatter payout
	var scatterPayout int
	if scatters > 0 {
		if val, ok := pays["B"][scatters]; ok {
			scatterPayout = val * spinReq.Bet / 100
		}
	}

	// line wins
	lineWins := s.EvaluateLines(board, spinReq)
	var lineTotal int
	for _, w := range lineWins {
		lineTotal += w.Payout
	}

	total := s.ApplyMaxPayout(lineTotal+scatterPayout, spinReq.Bet, maxPayoutMultiplier)

	awarded := 0
	if scatters >= 3 {
		if v, ok := s.cfg.FreeSpinsByScatter()[scatters]; ok {
			awarded = v

			countFreeSpins, err := s.lineRepo.GetCountFreeSpins(userID)
			if err != nil {
				return nil, errors.New("failed to get count free spins")
			}
			countFreeSpins += v
			err = s.lineRepo.UpdateCountFreeSpins(userID, countFreeSpins)
			if err != nil {
				return nil, errors.New("failed to update count free spins")
			}
		}
	}

	return &model.SpinResult{
		Board:            board,
		LineWins:         lineWins,
		ScatterCount:     scatters,
		ScatterPayout:    scatterPayout,
		AwardedFreeSpins: awarded,
		TotalPayout:      total,
	}, nil
}

// GenerateBoard генерирует игровое поле матрицы 5x3
func (s *LineService) GenerateBoard(forceBonus bool, spinReq model.LineSpin, userID string) ([5][3]string, error) {
	var board [5][3]string

	countFreeSpins, err := s.lineRepo.GetCountFreeSpins(userID)
	if err != nil {
		return board, errors.New("failed to get count free spins")
	}

	wildReels := map[int]bool{}
	for r := 1; r <= 3; r++ {
		if rand.Float64() < s.cfg.WildChance() || countFreeSpins > 0 {
			wildReels[r] = true
		}
	}

	scatterPositions := map[[2]int]bool{}
	if forceBonus || spinReq.BuyBonus {
		spinReq.BuyBonus = false
		cols := rand.Perm(reels)[:3]
		for _, reel := range cols {
			row := rand.Intn(rows)
			scatterPositions[[2]int{reel, row}] = true
		}
	}

	syms := s.cfg.SymbolWeights()
	for r := 0; r < reels; r++ {
		for row := 0; row < rows; row++ {
			if wildReels[r] {
				board[r][row] = "W"
			} else if scatterPositions[[2]int{r, row}] {
				board[r][row] = "B"
			} else {
				sym := s.RandomWeighted(syms)
				board[r][row] = sym
			}
		}
	}

	return board, nil
}

// EvaluateLines выполняет оценку выигрышных линий
func (s *LineService) EvaluateLines(board [5][3]string, spinReq model.LineSpin) []model.LineWin {
	var wins []model.LineWin
	for i, line := range playLines {
		symbols := make([]string, reels)
		for r := 0; r < 5; r++ {
			symbols[r] = board[r][line[r]]
		}
		if symbols[0] == "B" {
			continue
		}
		var base string
		for _, sym := range symbols {
			if sym != "W" && sym != "B" {
				base = sym
				break
			}
		}

		if base == "" {
			continue
		}
		count := 0
		for _, sym := range symbols {
			if sym == base || sym == "W" {
				count++
			} else {
				break
			}
		}
		if count >= 3 {
			if payTable, ok := pays[base]; ok {
				if val, ok := payTable[count]; ok {
					win := model.LineWin{
						Line:   i + 1,
						Symbol: base,
						Count:  count,
						Payout: val * spinReq.Bet / 100,
					}
					wins = append(wins, win)
				}
			}
		}
	}
	return wins
}

// RandomWeighted выполняет взвешенный случайный выбор символа
func (s *LineService) RandomWeighted(symbolWeights map[string]int) string {
	total := 0
	for _, w := range symbolWeights {
		total += w
	}
	if total <= 0 {
		for s := range symbolWeights {
			return s
		}
		return ""
	}
	r := rand.Intn(total)
	for s, w := range symbolWeights {
		if r < w {
			return s
		}
		r -= w
	}
	for s := range symbolWeights {
		return s
	}
	return ""
}

// ApplyMaxPayout применяет лимит по максимальному выигрышу
func (s *LineService) ApplyMaxPayout(amount, bet, maxMult int) int {
	maxPay := maxMult * bet
	if amount > maxPay {
		return maxPay
	}
	return amount
}
