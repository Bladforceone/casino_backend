package line_stats_repo

import (
	"casino_backend/internal/repository"
	"log"
	"sync"
)

const (
	// LimitHigh Средний avgPayoutPerBet > 1.5 → перейти к конфигу с более низкими шансами (увеличить индекс)
	LimitHigh = 1.5

	// LimitLow Средний avgPayoutPerBet < 0.8 → перейти к конфигу с более высокими шансами (уменьшить индекс)
	LimitLow = 0.8

	// MonitorWindow Количество спинов для мониторинга среднего выплат
	MonitorWindow = 100
)

type GameStats struct {
	History      []int64 // Массив с выплатами за спины (rolling window)
	TotalPayout  int64   // Сумма History
	SpinCount    int     // = len(History)
	CurrentIndex int     // 0-19, индекс конфига
}

type statsRepo struct {
	mtx  sync.RWMutex
	data GameStats
}

func NewLineStatsRepository() repository.LineStatsRepository {
	return &statsRepo{
		data: GameStats{
			CurrentIndex: 10, // Старт с середины (0-19)
			History:      make([]int64, 0),
		},
	}
}

// GetConfigIndex возвращает текущий индекс конфигурации.
func (r *statsRepo) GetConfigIndex() (int, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.data.CurrentIndex, nil
}

// UpdateStats обновляет статистику выплат и при необходимости переключает CurrentIndex.
func (r *statsRepo) UpdateStats(totalPayout int, bet int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	// Шаг 1: Обновить историю и метрики
	r.data.History = append(r.data.History, int64(totalPayout))
	r.data.TotalPayout += int64(totalPayout)
	if len(r.data.History) > MonitorWindow {
		oldPayout := r.data.History[0]
		r.data.History = r.data.History[1:]
		r.data.TotalPayout -= oldPayout
	}
	r.data.SpinCount = len(r.data.History)

	// Шаг 2: "Волшебство" — проверить пороги и переключить, если достаточно данных
	if r.data.SpinCount >= MonitorWindow {
		avgPayout := float64(r.data.TotalPayout) / float64(r.data.SpinCount)
		avgPayoutPerBet := avgPayout / float64(bet) // Нормализация на bet текущего спина
		if avgPayoutPerBet > LimitHigh {
			r.switchIndex(+1) // Уменьшить шансы (к низкому RTP)
		} else if avgPayoutPerBet < LimitLow {
			r.switchIndex(-1) // Увеличить шансы (к высокому RTP)
		}
	}
	return nil
}

// switchIndex изменяет CurrentIndex на +1 или -1 в зависимости от направления, и логирует изменение.
func (r *statsRepo) switchIndex(direction int) {
	r.data.CurrentIndex += direction
	if r.data.CurrentIndex < 0 {
		r.data.CurrentIndex = 0
	} else if r.data.CurrentIndex > 19 {
		r.data.CurrentIndex = 19
	}
	// Логируем смену индекса для отладки
	log.Printf("[Line] Switched to config index %d", r.data.CurrentIndex)
}
