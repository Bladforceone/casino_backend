package cascade

import (
	"casino_backend/internal/config"
	"casino_backend/internal/middleware"
	"context"
	"errors"
	"math/rand"

	"casino_backend/internal/model"
)

const (
	rows = 7
	cols = 7

	// Символы: 0..6 обычные (по возрастанию ценности), 7 - бонусный (scatter-like в логике сбора)
	//symbolRegularCount = 7
	symbolBonus = 7

	// Множители на ячейках: при втором удалении x2, далее удваивается до MAX
	multiplierStart = 2
	multiplierMax   = 128 // До x128 — чтобы не переполнить int при умножении

	// Предел итераций разрешения каскадов
	maxResolveIter = 100

	// Ограничение максимального выигрыша (в кратности ставки)
	maxWinXBet = 10000
)

// Пустая ячейка
const emptyCell = -1

type cluster struct {
	symbol int
	cells  [][2]int
}

// Spin — основной метод
func (s *serv) Spin(ctx context.Context, req model.CascadeSpin) (*model.CascadeSpinResult, error) {
	// Валидация ставки
	// Если ставка меньше либо равна нулю или не кратна 2-м (т.е. Нечетная) — ошибка
	if req.Bet <= 0 || req.Bet%2 != 0 {
		return nil, errors.New("bet must be positive and even")
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.New("user id not found in context")
	}

	// Получаем текущий индекс конфига из статистики (вне транзакции)
	configIndex, err := s.cascadeStatsRepo.GetConfigIndex()
	if err != nil {
		return nil, err
	}
	// Выбираем конфиг по индексу
	currentCfg := s.cfg[configIndex]

	var spinRes *model.CascadeSpinResult
	var finalFreeSpins int

	// Начало транзакции
	err = s.txManager.Do(ctx, func(txCtx context.Context) error {
		// Проверка баланса и фриспинов внутри транзакции
		freeSpins, err := s.cascadeRepo.GetFreeSpinCount(txCtx, userID)
		if err != nil {
			return err
		}

		isFreeSpin := freeSpins > 0
		var userBalance int

		if !isFreeSpin {
			userBalance, err = s.userRepo.GetBalance(txCtx, userID)
			if err != nil {
				return err
			}
			if userBalance < req.Bet {
				return errors.New("not enough balance")
			}
			userBalance -= req.Bet
			if err := s.userRepo.UpdateBalance(txCtx, userID, userBalance); err != nil {
				return err
			}
		} else {
			freeSpins--
			if err := s.cascadeRepo.UpdateFreeSpinCount(txCtx, userID, freeSpins); err != nil {
				return err
			}
			// Получаем баланс для последующего начисления
			userBalance, err = s.userRepo.GetBalance(txCtx, userID)
			if err != nil {
				return err
			}
		}

		// Выполняем спин (с txCtx)
		spinRes, err = s.spinOnce(txCtx, userID, req.Bet, !isFreeSpin, currentCfg)
		if err != nil {
			return err
		}

		// Начисление выигрыша
		userBalance += spinRes.TotalPayout
		if err := s.userRepo.UpdateBalance(txCtx, userID, userBalance); err != nil {
			return err
		}

		// Начисление фриспинов — уже сделано внутри spinOnce → просто читаем результат
		finalFreeSpins, err = s.cascadeRepo.GetFreeSpinCount(txCtx, userID)
		if err != nil {
			return err
		}
		spinRes.FreeSpinsLeft = finalFreeSpins

		// Заполняем индексы каскадов (0 = первый)
		for i := range spinRes.Cascades {
			spinRes.Cascades[i].CascadeIndex = i
		}

		// Сохраняем balance для возврата
		spinRes.Balance = userBalance
		spinRes.InFreeSpin = isFreeSpin

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Обновляем статистику (вне транзакции)
	err = s.cascadeStatsRepo.UpdateStats(spinRes.TotalPayout, req.Bet)
	if err != nil {
		return nil, errors.New("failed to update stats")
	}

	return &model.CascadeSpinResult{
		InitialBoard:     spinRes.InitialBoard,
		Board:            spinRes.Board,
		Cascades:         spinRes.Cascades,
		TotalPayout:      spinRes.TotalPayout,
		Balance:          spinRes.Balance,
		ScatterCount:     spinRes.ScatterCount,
		AwardedFreeSpins: spinRes.AwardedFreeSpins,
		FreeSpinsLeft:    finalFreeSpins,
		InFreeSpin:       spinRes.InFreeSpin,
	}, nil
}

// spinOnce полный спин с каскадами
func (s *serv) spinOnce(ctx context.Context, userID int, bet int, resetMultipliers bool, cfg config.CascadeConfig) (*model.CascadeSpinResult, error) {
	// Инициализация доски
	var board [rows][cols]int
	// hits - сколько раз ячейка участвовала в удалении кластера
	// mult - множитель клетки (x1, x2, x4, x8, x16...)
	// Загружаем состояние множителей из репозитория
	mult, hits, err := s.cascadeRepo.GetMultiplierState(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Если обычный спин, то сбрасываем множители и заново их инициализируем
	if resetMultipliers {
		// Обнуляем счетчики и множители в репозитории
		if err := s.cascadeRepo.ResetMultiplierState(ctx, userID); err != nil {
			return nil, err
		}
		// Загружаем заново — Reset уже поставил 1 и 0
		mult, hits, err = s.cascadeRepo.GetMultiplierState(ctx, userID)
		if err != nil {
			return nil, err
		}
		// Заполняем доску заново
		s.fillBoard(&board, cfg.BonusProbPerColumn(), cfg.SymbolWeights())
	} else {
		// Фриспин — оставляем старые множители, но генерим новую доску
		s.fillBoard(&board, cfg.BonusProbPerColumn(), cfg.SymbolWeights())
		// ← Важно: множители остаются от прошлого спина!
	}

	// Сохраняем начальную доску для возврата
	initialBoard := board

	// Инициализируем каскады
	var cascades []model.CascadeStep
	// Общий выигрыш за спин
	var totalWin int

	for iter := 0; iter < maxResolveIter; iter++ {
		clusters := s.findClusters(board)
		if len(clusters) == 0 {
			break
		}

		step := model.CascadeStep{}

		// Обрабатываем все кластеры на доске (подсчет выигрыша, удаление, обновление множителей)
		for _, cl := range clusters {
			win := s.calculateWin(cl, mult, bet, cfg.PayoutTable())
			totalWin += win
			avgMult := s.averageMultiplier(cl, mult)

			positions := make([]model.Position, len(cl.cells))
			for i, cell := range cl.cells {
				positions[i] = model.Position{Row: cell[0], Col: cell[1]}
			}

			step.Clusters = append(step.Clusters, model.ClusterInfo{
				Symbol:     cl.symbol,
				Cells:      positions,
				Count:      len(cl.cells),
				Payout:     win,
				Multiplier: avgMult,
			})

			s.removeCluster(cl, &board, &hits, &mult)
		}

		// Сдвигаем символы вниз и заполняем пустоты
		s.collapse(&board)
		intermediateBoard := board // Копия после collapse (upper empty)
		s.refill(&board, cfg.BonusProbPerColumn(), cfg.SymbolWeights())

		// Добавляем новые символы которые упадут на доску
		step.NewSymbols = []struct {
			model.Position
			Symbol int
		}{}
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				if intermediateBoard[r][c] == emptyCell && board[r][c] != emptyCell {
					step.NewSymbols = append(step.NewSymbols, struct {
						model.Position
						Symbol int
					}{Position: model.Position{Row: r, Col: c}, Symbol: board[r][c]})
				}
			}
		}
		cascades = append(cascades, step)
	}

	// Сохраняем обновлённое состояние множителей
	if err := s.cascadeRepo.SetMultiplierState(ctx, userID, mult, hits); err != nil {
		return nil, err
	}

	scatterCount := s.countScatters(board)
	awarded := 0
	if scatterCount >= 3 {
		if v, ok := cfg.BonusAwards()[scatterCount]; ok {
			awarded = v

			currentFS, err := s.cascadeRepo.GetFreeSpinCount(ctx, userID)
			if err != nil {
				return nil, err
			}
			err = s.cascadeRepo.UpdateFreeSpinCount(ctx, userID, currentFS+awarded)
			if err != nil {
				return nil, err
			}
		}
	}
	totalPayout := s.applyMaxPayout(totalWin, bet)

	return &model.CascadeSpinResult{
		InitialBoard:     initialBoard,
		Board:            board,
		Cascades:         cascades,
		TotalPayout:      totalPayout,
		ScatterCount:     scatterCount,
		AwardedFreeSpins: awarded,
	}, nil
}

//---------- ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ ----------

// fillBoard заполняет доску начальными символами
func (s *serv) fillBoard(board *[rows][cols]int, bonusProbPerColumn float64, weights map[int]int) {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if rand.Float64() < bonusProbPerColumn {
				board[r][c] = symbolBonus
			} else {
				board[r][c] = s.randomRegularSymbol(weights)
			}
		}
	}
}

// collapse сдвигает символы вниз, устанавливает upper empty
func (s *serv) collapse(board *[rows][cols]int) {
	for c := 0; c < cols; c++ {
		stack := make([]int, 0, rows)
		for r := 0; r < rows; r++ {
			if board[r][c] != emptyCell {
				stack = append(stack, board[r][c])
			}
		}
		for r := 0; r < rows; r++ { // Сначала очистим всю колонку
			board[r][c] = emptyCell
		}
		for i, sym := range stack {
			board[rows-len(stack)+i][c] = sym // Сдвиг вниз (bottom)
		}
		// Upper уже empty
	}
}

// refill заполняет empty (upper) новыми символами
func (s *serv) refill(board *[rows][cols]int, bonusProbPerColumn float64, weights map[int]int) {
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			if board[r][c] == emptyCell {
				if rand.Float64() < bonusProbPerColumn {
					board[r][c] = symbolBonus
				} else {
					board[r][c] = s.randomRegularSymbol(weights)
				}
			}
		}
	}
}

// randomRegularSymbol выбирает случайный обычный символ с учётом весов
func (s *serv) randomRegularSymbol(weights map[int]int) int {

	total := 0
	for _, w := range weights {
		total += w
	}
	if total == 0 {
		return 0
	}
	n := rand.Intn(total)
	for sym, w := range weights {
		if n < w {
			return sym
		}
		n -= w
	}
	return 0
}

// findClusters ищет кластеры на доске
func (s *serv) findClusters(board [rows][cols]int) []cluster {
	visited := [rows][cols]bool{}
	var clusters []cluster
	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if visited[r][c] || board[r][c] == emptyCell || board[r][c] == symbolBonus {
				continue
			}
			sym := board[r][c]
			var component [][2]int
			queue := [][2]int{{r, c}}
			visited[r][c] = true

			for len(queue) > 0 {
				cur := queue[0]
				queue = queue[1:]
				component = append(component, cur)
				cr, cc := cur[0], cur[1]
				for _, d := range dirs {
					nr, nc := cr+d[0], cc+d[1]
					if nr >= 0 && nr < rows && nc >= 0 && nc < cols &&
						!visited[nr][nc] && board[nr][nc] == sym {
						visited[nr][nc] = true
						queue = append(queue, [2]int{nr, nc})
					}
				}
			}
			if len(component) >= 5 {
				clusters = append(clusters, cluster{symbol: sym, cells: component})
			}
		}
	}
	return clusters
}

// calculateWin вычисляет выигрыш за кластер
func (s *serv) calculateWin(cl cluster, mult [rows][cols]int, bet int, payTable map[int]int) int {
	// Защита от пустого кластера (на всякий случай, хотя findClusters фильтрует >=5)
	length := len(cl.cells)
	if length == 0 {
		return 0
	}

	base, ok := payTable[cl.symbol]
	if !ok {
		base = 0 // или можно логгировать ошибку конфигурации
	}

	// Базовая выплата: base × количество символов
	baseWin := base * length

	// Суммируем множители по всем ячейкам кластера
	var sumMult int
	for _, cell := range cl.cells {
		sumMult += mult[cell[0]][cell[1]]
	}

	// Средний множитель (округление вниз — как в оригинале)
	avgMult := sumMult / length
	if avgMult < 1 {
		avgMult = 1
	}

	return baseWin * avgMult * bet
}

// averageMultiplier возвращает средний множитель кластера (для отображения клиенту)
func (s *serv) averageMultiplier(cl cluster, mult [rows][cols]int) int {
	length := len(cl.cells)
	if length == 0 {
		return 1
	}

	var sum int
	for _, cell := range cl.cells {
		sum += mult[cell[0]][cell[1]]
	}

	avg := sum / length
	if avg < 1 {
		avg = 1
	}
	return avg
}

// removeCluster удаляет кластер с доски и обновляет счётчики попаданий и множители
func (s *serv) removeCluster(cl cluster, board *[rows][cols]int, hits *[rows][cols]int, mult *[rows][cols]int) {
	for _, cell := range cl.cells {
		r, c := cell[0], cell[1]
		hits[r][c]++
		if hits[r][c] >= 2 {
			shift := uint(hits[r][c] - 2)
			newMult := multiplierStart << shift
			if newMult > multiplierMax {
				newMult = multiplierMax
			}
			mult[r][c] = newMult
		}
		board[r][c] = emptyCell
	}
}

// countScatters подсчитывает количество бонусных символов на доске
func (s *serv) countScatters(board [rows][cols]int) int {
	cnt := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if board[r][c] == symbolBonus {
				cnt++
			}
		}
	}
	return cnt
}

// applyMaxPayout ограничивает максимальный выигрыш
func (s *serv) applyMaxPayout(amount, bet int) int {
	maxAllowed := maxWinXBet * bet
	if amount > maxAllowed {
		return maxAllowed
	}
	return amount
}
