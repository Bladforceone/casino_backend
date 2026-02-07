package line

import (
	"casino_backend/internal/middleware"
	"casino_backend/internal/model"
	servModel "casino_backend/internal/service/line/model"
	"context"
	"errors"
	"log"
	"math/rand"
)

const (
	// Барабаны
	reels = 5
	// Линии
	rows = 3
	// Максимальная выплата в кратности ставки
	maxPayoutMultiplier = 10000
)

// Spin выполняет спин с учётом баланса и фриспинов
func (s *serv) Spin(ctx context.Context, spinReq model.LineSpin) (*model.SpinResult, error) {
	// Валидация ставки
	// Если ставка меньше либо равна нулю или не кратна 2-м (т.е. Нечетная) — ошибка
	if spinReq.Bet <= 0 || spinReq.Bet%2 != 0 {
		return nil, errors.New("bet must be positive and even")
	}

	// Получаем ID пользователя
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.New("user id not found in context")
	}

	// Получаем пресет весов символов исходя из статистики
	presetCfg := servModel.RtpPresets[s.lineStatsRepo.CasinoState().PresetIndex]

	// Инициализируем структуру для хранения результатов спина
	var res *model.SpinResult

	// Начало транзакции где выполняется процесс спина.
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		// Получаем текущее количество фриспинов внутри транзакции
		countFreeSpins, err := s.repo.GetFreeSpinCount(txCtx, userID)
		if err != nil {
			// Елси этих данных нет, то значит создаем их по умолчанию
			err = s.repo.CreateLineGameState(ctx, userID)
			if err != nil {
				log.Println(err)
				return errors.New("failed to get count free spins in Line Repo")
			}
			countFreeSpins = 0
		}

		// Локальная переменная для баланса
		var userBalance int

		//TODO Бойлерплейт ниже. Два раза получаем баланс когда можно получить один раз до условия

		// Платный спин
		// Если счетчик фриспинов нулевой, то списываем деньги с баланса
		if countFreeSpins == 0 {
			// Получаем баланс пользователя
			userBalance, err = s.userRepo.GetBalance(txCtx, userID)
			if err != nil {
				return errors.New("failed to get user balance")
			}
			if userBalance < spinReq.Bet {
				return errors.New("not enough balance")
			}

			// Списание ставки, обновление баланса пользователя
			userBalance -= spinReq.Bet
			if err := s.userRepo.UpdateBalance(txCtx, userID, userBalance); err != nil {
				return errors.New("failed to update user balance")
			}
		} else { // Иначе режим фриспинов.
			// Уменьшаем счетчик фриспинов на 1
			if err := s.repo.UpdateFreeSpinCount(txCtx, userID, countFreeSpins-1); err != nil {
				return errors.New("failed to update count free spins")
			}
			// Получаем баланс для последующего начисления
			userBalance, err = s.userRepo.GetBalance(txCtx, userID)
			if err != nil {
				return errors.New("failed to get user balance")
			}
		}

		// КЛЮЧЕВОЙ ВЫЗОВ
		// Делаем спин (передаём countFreeSpins как параметр)
		res, err = s.SpinOnce(spinReq, presetCfg, countFreeSpins)
		if err != nil {
			return err
		}

		// Устанавливаем флаг InFreeSpin, если это был фриспин
		if countFreeSpins > 0 {
			res.InFreeSpin = true
		}

		// Начисление выигрыша
		userBalance += res.TotalPayout
		if err := s.userRepo.UpdateBalance(txCtx, userID, userBalance); err != nil {
			return errors.New("failed to update user balance")
		}

		// Если есть выигранные фриспины, добавляем их
		if res.AwardedFreeSpins > 0 {
			// Получаем текущее количество фриспинов (после возможного уменьшения)
			currentFree, err := s.repo.GetFreeSpinCount(txCtx, userID)
			if err != nil {
				return errors.New("failed to get count free spins")
			}
			// Прибавляем новые спины
			if err := s.repo.UpdateFreeSpinCount(txCtx, userID, currentFree+res.AwardedFreeSpins); err != nil {
				return errors.New("failed to update count free spins")
			}
			// Обновляем то, что увидит клиент
			res.FreeSpinCount = currentFree + res.AwardedFreeSpins
		}

		// Получаем финальное количество фриспинов для возврата
		freeCount, err := s.repo.GetFreeSpinCount(txCtx, userID)
		if err != nil {
			return errors.New("failed to get count free spins")
		}

		// Устанавливаем финальные значения в res
		res.Balance = userBalance
		res.FreeSpinCount = freeCount // Финальное значение (перезапишет, если было awarded)

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Обновляем статистику
	err = s.lineStatsRepo.UpdateState(float64(spinReq.Bet), float64(res.TotalPayout))
	if err != nil {
		return nil, errors.New("failed to update stats")
	}
	return res, nil
}

// SpinOnce выполняет один спин (возвращает единый SpinResult)
func (s *serv) SpinOnce(spinReq model.LineSpin, preset servModel.RTPPreset, countFreeSpins int) (*model.SpinResult, error) {
	board, err := s.GenerateBoard(countFreeSpins, cfg.WildChance(idx), cfg.SymbolWeights(idx))
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
		if val, ok := servModel.PayoutTable["B"][scatters]; ok {
			scatterPayout = val * spinReq.Bet / 100
		}
	}

	// line wins
	lineWins := s.EvaluateLines(board, spinReq, servModel.PayoutTable)
	var lineTotal int
	for _, w := range lineWins {
		lineTotal += w.Payout
	}

	total := s.ApplyMaxPayout(lineTotal+scatterPayout, spinReq.Bet, maxPayoutMultiplier)

	awarded := 0
	if scatters >= 3 {
		if v, ok := servModel.FreeSpinsScatter[scatters]; ok {
			awarded = v
		}
	}

	return &model.SpinResult{
		Board:            board,
		LineWins:         lineWins,
		ScatterCount:     scatters,
		ScatterPayout:    scatterPayout,
		AwardedFreeSpins: awarded,
		TotalPayout:      total,
		Balance:          0,
	}, nil
}

// GenerateBoard генерирует игровое поле матрицы 5x3
func (s *serv) GenerateBoard(countFreeSpins int, wildChance float64, symbolWeights map[string]int) ([5][3]string, error) {
	var board [5][3]string

	// Добавляем вайлды только на центральные 3 барабана (индексы 1,2,3)
	wildReels := map[int]bool{}
	if countFreeSpins > 0 {
		// ГАРАНТИРОВАННО хотя бы один Wild каждый спин бонуски
		guaranteedReel := 1 + rand.Intn(3) // 1, 2 или 3 → барабаны 2,3,4
		wildReels[guaranteedReel] = true
		// Остальные два барабана могут тоже стать Wild с шансом 6%
		for reel := 1; reel <= 3; reel++ {
			if reel != guaranteedReel && rand.Float64() < wildChance {
				wildReels[reel] = true
			}
		}
	} else { // Обычная игра — обычный шанс 6% на каждый центральный барабан
		for reel := 1; reel <= 3; reel++ {
			if rand.Float64() < wildChance {
				wildReels[reel] = true
			}
		}
	}

	// Отслеживаем, уже выпал ли скаттер на этом барабане
	hasScatter := make([]bool, reels) // false по умолчанию

	// Заполняем остальное случайными символами
	for r := 0; r < reels; r++ {
		for row := 0; row < rows; row++ {
			if wildReels[r] {
				board[r][row] = "W"
				continue
			}

			var sym string
			if hasScatter[r] {
				// На этом барабане уже есть скаттер → больше нельзя
				sym = s.RandomWeightedNoScatter(symbolWeights)
			} else {
				// Обычный ролл, скаттер ещё разрешён
				sym = s.RandomWeighted(symbolWeights)
			}

			board[r][row] = sym

			// Если только что выпал скаттер — помечаем барабан
			if sym == "B" {
				hasScatter[r] = true
			}
		}
	}
	return board, nil
}

// EvaluateLines выполняет оценку выигрышных линий
func (s *serv) EvaluateLines(board [5][3]string, spinReq model.LineSpin, payoutTable map[string]map[int]int) []model.LineWin {
	var wins []model.LineWin
	for i, line := range servModel.PlayLines {
		symbols := make([]string, reels)
		for r := 0; r < 5; r++ {
			symbols[r] = board[r][line[r]]
		}

		// Пропускаем линии, где первый символ — скаттер
		if symbols[0] == "B" {
			continue
		}

		// Находим базовый символ (не W и не B)
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

		// Считаем последовательность base + W с первого барабана
		count := 0
		for _, sym := range symbols {
			if sym == base || sym == "W" {
				count++
			} else {
				break
			}
		}

		// Определяем минимальное количество символов для выплаты
		minCount := 3
		for c := range payoutTable[base] {
			if c < minCount {
				minCount = c // обновится до 2 для S8
			}
		}

		if count >= minCount {
			if payTable, ok := payoutTable[base]; ok {
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
func (s *serv) RandomWeighted(symbolWeights map[string]int) string {
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

// RandomWeightedNoScatter — выбирает символ по весам, но полностью исключает скаттер "B"
func (s *serv) RandomWeightedNoScatter(symbolWeights map[string]int) string {
	total := 0
	for sym, w := range symbolWeights {
		if sym != "B" { // полностью игнорируем скаттер
			total += w
		}
	}
	if total <= 0 {
		for sym := range symbolWeights {
			if sym != "B" {
				return sym
			}
		}
		return ""
	}
	r := rand.Intn(total)
	current := 0
	for sym, w := range symbolWeights {
		if sym == "B" {
			continue
		}
		if r < current+w {
			return sym
		}
		current += w
	}

	// fallback
	for sym := range symbolWeights {
		if sym != "B" {
			return sym
		}
	}
	return ""
}

// ApplyMaxPayout применяет лимит по максимальному выигрышу
func (s *serv) ApplyMaxPayout(amount, bet, maxMult int) int {
	maxPay := maxMult * bet
	if amount > maxPay {
		return maxPay
	}
	return amount
}
