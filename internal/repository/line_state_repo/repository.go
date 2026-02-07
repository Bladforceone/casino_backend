package line_state_repo

import (
	repoModel "casino_backend/internal/repository/line_state_repo/model"
	"sync"
)

// Реализация репозитория для хранения состояния казино
type stateRepo struct {
	mtx   sync.RWMutex
	state repoModel.CasinoState
}

// Конструктор для создания нового репозитория с начальным состоянием
func NewLineStatsRepository(initialState repoModel.CasinoState) *stateRepo {
	return &stateRepo{
		state: initialState,
	}
}

// Получение текущего состояния казино
// Является геттером для структуры состояния казино
// Возвращает копию структуры CasinoState
func (r *stateRepo) CasinoState() repoModel.CasinoState {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.state
}

// Обновление состояния казино после спина
func (r *stateRepo) UpdateState(bet, payout float64) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.state.TotalSpins++
	r.state.TotalBet += bet
	r.state.TotalPayout += payout
	if r.state.TotalBet > 0 {
		r.state.CurrentRTP = r.state.TotalPayout / r.state.TotalBet * 100
	}

	// Добавляем спин в окно
	spinRTP := 0.0
	if bet > 0 {
		spinRTP = payout / bet * 100
	}
	r.state.SpinWindow = append(r.state.SpinWindow, repoModel.SpinResult{
		Bet:    bet,
		Payout: payout,
		RTP:    spinRTP,
	})

	// Поддерживаем размер окна
	if len(r.state.SpinWindow) > r.state.WindowSize {
		r.state.SpinWindow = r.state.SpinWindow[1:]
	}

	// Пересчитываем RTP в окне (объеденил с функцией recalculateWindowRTP)
	var windowBet, windowPayout float64
	for _, spin := range r.state.SpinWindow {
		windowBet += spin.Bet
		windowPayout += spin.Payout
	}

	if windowBet > 0 {
		r.state.WindowRTP = windowPayout / windowBet * 100
	} else {
		r.state.WindowRTP = 0
	}
}
