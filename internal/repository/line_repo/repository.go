package line_repo

import (
	"casino_backend/internal/repository"
	"sync"
)

type memoryData struct {
	balance       int
	freeSpinCount int
}

type repo struct {
	mtx sync.RWMutex
	mem memoryData
}

func NewLineRepository() repository.LineRepository {
	return &repo{mem: memoryData{}}
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
