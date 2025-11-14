package tests

import (
	"casino_backend/internal/model"
	"casino_backend/internal/service/line"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// TestLineService_Spin_PaidSpinSuccess тестирует успешный платный спин без фриспинов
func TestLineService_Spin_PaidSpinSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	lineRepo := NewMockLineRepository(ctrl)
	cfg := NewMockLineConfig(ctrl)

	service := line.NewLineService(cfg, userRepo, lineRepo)
	ctx := context.WithValue(context.Background(), "userID", "user123")

	spinReq := model.LineSpin{
		Bet:      100,
		BuyBonus: false,
	}

	// Setup expectations
	lineRepo.ExpectGetCountFreeSpins("user123", 0, nil) // First: в Spin() начало
	userRepo.ExpectGetBalance("user123", 1000, nil)     // платный спин
	userRepo.ExpectUpdateBalance("user123", 900, nil)   // списание ставки
	lineRepo.ExpectGetCountFreeSpins("user123", 0, nil) // Не будет скаттеров, скорее всего не вызовется
	userRepo.ExpectGetBalance("user123", 900, nil)      // обновление баланса
	userRepo.ExpectUpdateBalance("user123", 1100, nil)  // добавляем выигрыш
	lineRepo.ExpectGetCountFreeSpins("user123", 0, nil) // Last: в Spin() конец

	cfg.ExpectSymbolWeights(map[string]int{
		"S1": 10, "S2": 10, "S3": 10, "S4": 10, "S5": 10,
		"S6": 5, "S7": 5, "S8": 2,
	})
	cfg.ExpectWildChance(0.1)
	cfg.ExpectFreeSpinsByScatter(map[int]int{
		3: 5, 4: 10, 5: 15,
	})

	result, err := service.Spin(ctx, spinReq)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.InFreeSpin)
	require.Equal(t, 0, result.FreeSpinCount)
}

// TestLineService_Spin_InsufficientBalance тестирует платный спин без достаточного баланса
func TestLineService_Spin_InsufficientBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	lineRepo := NewMockLineRepository(ctrl)
	cfg := NewMockLineConfig(ctrl)

	service := line.NewLineService(cfg, userRepo, lineRepo)
	ctx := context.WithValue(context.Background(), "userID", "user123")

	spinReq := model.LineSpin{
		Bet:      100,
		BuyBonus: false,
	}

	lineRepo.ExpectGetCountFreeSpins("user123", 0, nil)
	userRepo.ExpectGetBalance("user123", 50, nil)

	result, err := service.Spin(ctx, spinReq)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "not enough balance", err.Error())
}

// TestLineService_Spin_FreeSpinSuccess тестирует успешный фриспин
func TestLineService_Spin_FreeSpinSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	lineRepo := NewMockLineRepository(ctrl)
	cfg := NewMockLineConfig(ctrl)

	service := line.NewLineService(cfg, userRepo, lineRepo)
	ctx := context.WithValue(context.Background(), "userID", "user123")

	spinReq := model.LineSpin{
		Bet:      100,
		BuyBonus: false,
	}

	lineRepo.ExpectGetCountFreeSpins("user123", 5, nil)    // First: в Spin() начало
	lineRepo.ExpectUpdateCountFreeSpins("user123", 4, nil) // фриспин - уменьшение
	lineRepo.ExpectGetCountFreeSpins("user123", 4, nil)    // SpinOnce - возможно, если будут скаттеры
	userRepo.ExpectGetBalance("user123", 1000, nil)        // обновление баланса
	userRepo.ExpectUpdateBalance("user123", 1200, nil)     // добавляем выигрыш
	lineRepo.ExpectGetCountFreeSpins("user123", 4, nil)    // Last: в Spin() конец

	cfg.ExpectSymbolWeights(map[string]int{
		"S1": 10, "S2": 10, "S3": 10, "S4": 10, "S5": 10,
		"S6": 5, "S7": 5, "S8": 2,
	})
	cfg.ExpectWildChance(0.1)
	cfg.ExpectFreeSpinsByScatter(map[int]int{
		3: 5, 4: 10, 5: 15,
	})

	result, err := service.Spin(ctx, spinReq)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.InFreeSpin)
	require.Equal(t, 4, result.FreeSpinCount)
}

// TestLineService_Spin_NoUserID тестирует отсутствие userID в контексте
func TestLineService_Spin_NoUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	lineRepo := NewMockLineRepository(ctrl)
	cfg := NewMockLineConfig(ctrl)

	service := line.NewLineService(cfg, userRepo, lineRepo)
	ctx := context.Background() // Нет userID

	spinReq := model.LineSpin{
		Bet:      100,
		BuyBonus: false,
	}

	result, err := service.Spin(ctx, spinReq)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "user ID not found in context", err.Error())
}

// TestLineService_Spin_GetBalanceError тестирует ошибку получения баланса
func TestLineService_Spin_GetBalanceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	lineRepo := NewMockLineRepository(ctrl)
	cfg := NewMockLineConfig(ctrl)

	service := line.NewLineService(cfg, userRepo, lineRepo)
	ctx := context.WithValue(context.Background(), "userID", "user123")

	spinReq := model.LineSpin{
		Bet:      100,
		BuyBonus: false,
	}

	lineRepo.ExpectGetCountFreeSpins("user123", 0, nil)
	userRepo.ExpectGetBalance("user123", 0, errors.New("db connection error"))

	result, err := service.Spin(ctx, spinReq)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "failed to get user balance", err.Error())
}

// TestLineService_Spin_UpdateBalanceError тестирует ошибку обновления баланса
func TestLineService_Spin_UpdateBalanceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	lineRepo := NewMockLineRepository(ctrl)
	cfg := NewMockLineConfig(ctrl)

	service := line.NewLineService(cfg, userRepo, lineRepo)
	ctx := context.WithValue(context.Background(), "userID", "user123")

	spinReq := model.LineSpin{
		Bet:      100,
		BuyBonus: false,
	}

	lineRepo.ExpectGetCountFreeSpins("user123", 0, nil)
	userRepo.ExpectGetBalance("user123", 1000, nil)
	userRepo.ExpectUpdateBalance("user123", 900, errors.New("db connection error"))

	result, err := service.Spin(ctx, spinReq)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "failed to update user balance", err.Error())
}

// TestLineService_EvaluateLines тестирует оценку выигрышных линий
func TestLineService_EvaluateLines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := NewMockLineConfig(ctrl)
	service := line.NewLineService(cfg, nil, nil)

	spinReq := model.LineSpin{Bet: 100}

	// Тестовая доска: первая линия (playLines[0] = {1, 1, 1, 1, 1}) должна быть S1 на всех позициях
	// board[reel][row], поэтому для row=1: board[r][1] должны быть все S1
	board := [5][3]string{
		{"XX", "S1", "S2"}, // reel 0: row 0=XX, row 1=S1, row 2=S2
		{"XX", "S1", "S3"}, // reel 1: row 0=XX, row 1=S1, row 2=S3
		{"XX", "S1", "S4"}, // reel 2
		{"XX", "S1", "S5"}, // reel 3
		{"XX", "S1", "S6"}, // reel 4
	}

	wins := service.EvaluateLines(board, spinReq)

	// Проверяем, что нашли выигрыш на первой линии (все S1)
	require.Greater(t, len(wins), 0)

	foundFirstLine := false
	for _, win := range wins {
		if win.Line == 1 && win.Symbol == "S1" && win.Count == 5 {
			foundFirstLine = true
			require.Equal(t, 450, win.Payout) // S1: 5 = 450 * 100 / 100 = 450
		}
	}
	require.True(t, foundFirstLine)
}

// TestLineService_EvaluateLines_WithWild тестирует оценку линий с дикой картой
func TestLineService_EvaluateLines_WithWild(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := NewMockLineConfig(ctrl)
	service := line.NewLineService(cfg, nil, nil)

	spinReq := model.LineSpin{Bet: 100}

	// Линия: S1 S1 W S1 S1 (где W - wild/дикая карта) на line {1, 1, 1, 1, 1}
	board := [5][3]string{
		{"XX", "S1", "S2"}, // reel 0: row 1 = S1
		{"XX", "S1", "S3"}, // reel 1: row 1 = S1
		{"XX", "W", "S4"},  // reel 2: row 1 = W (wild)
		{"XX", "S1", "S5"}, // reel 3: row 1 = S1
		{"XX", "S1", "S6"}, // reel 4: row 1 = S1
	}

	wins := service.EvaluateLines(board, spinReq)

	// Первая линия должна иметь 5 совпадений (S1 с дикими картами)
	foundLine := false
	for _, win := range wins {
		if win.Line == 1 && win.Symbol == "S1" && win.Count == 5 {
			foundLine = true
		}
	}
	require.True(t, foundLine)
}

// TestLineService_EvaluateLines_ScatterBreak тестирует, что линии не начинаются со scatter
func TestLineService_EvaluateLines_ScatterBreak(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := NewMockLineConfig(ctrl)
	service := line.NewLineService(cfg, nil, nil)

	spinReq := model.LineSpin{Bet: 100}

	// Первая линия начинается со scatter (B)
	board := [5][3]string{
		{"B", "S5", "S2"},
		{"S1", "S3", "S4"},
		{"S1", "S6", "S7"},
		{"S1", "S7", "S8"},
		{"S1", "S2", "W"},
	}

	wins := service.EvaluateLines(board, spinReq)

	// Первая линия не должна быть в результатах
	for _, win := range wins {
		require.NotEqual(t, 1, win.Line)
	}
}

// TestLineService_RandomWeighted тестирует взвешенный выбор символов
func TestLineService_RandomWeighted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := NewMockLineConfig(ctrl)
	service := line.NewLineService(cfg, nil, nil)

	weights := map[string]int{
		"A": 100,
		"B": 1,
	}

	// Много итераций для статистической проверки
	counts := make(map[string]int)
	for i := 0; i < 10000; i++ {
		sym := service.RandomWeighted(weights)
		counts[sym]++
	}

	// A должен появляться гораздо чаще, чем B
	require.Greater(t, counts["A"], counts["B"])
}

// TestLineService_ApplyMaxPayout тестирует применение лимита выплаты
func TestLineService_ApplyMaxPayout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := NewMockLineConfig(ctrl)
	service := line.NewLineService(cfg, nil, nil)

	tests := []struct {
		name     string
		amount   int
		bet      int
		maxMult  int
		expected int
	}{
		{"Обычная выплата", 500, 100, 10000, 500},
		{"Выплата ниже лимита", 1000000, 100, 10000, 1000000},
		{"Выплата выше лимита", 2000000, 100, 10000, 1000000},
		{"Граничный случай", 1000000, 100, 10000, 1000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ApplyMaxPayout(tt.amount, tt.bet, tt.maxMult)
			require.Equal(t, tt.expected, result)
		})
	}
}
