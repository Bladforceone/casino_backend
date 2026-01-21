package cascade_repo

import (
	"casino_backend/internal/repository"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type memoryData struct {
	balance       int
	freeSpinCount int
	mult          [7][7]int // Множители
	hits          [7][7]int // Счётчики попаданий
}

type repo struct {
	mtx sync.RWMutex
	mem memoryData
	dbc *pgxpool.Pool
}

func NewCascadeRepository(dbc *pgxpool.Pool) repository.CascadeRepository {
	return &repo{
		mem: memoryData{},
		dbc: dbc,
	}
}

func (r *repo) GetBalance() (int, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.mem.balance, nil
}

func (r *repo) UpdateBalance(amount int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.mem.balance = amount
	return nil
}

func (r *repo) GetFreeSpinCount() (int, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.mem.freeSpinCount, nil
}

func (r *repo) UpdateFreeSpinCount(count int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.mem.freeSpinCount = count
	return nil
}

// ResetMultiplierState Сброс при начале платного спина
func (r *repo) ResetMultiplierState() error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for i := range r.mem.mult {
		for j := range r.mem.mult[i] {
			r.mem.mult[i][j] = 1
			r.mem.hits[i][j] = 0
		}
	}
	return nil
}

func (r *repo) GetMultiplierState() ([7][7]int, [7][7]int) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.mem.mult, r.mem.hits
}

func (r *repo) SetMultiplierState(mult, hits [7][7]int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.mem.mult = mult
	r.mem.hits = hits
	return nil
}
