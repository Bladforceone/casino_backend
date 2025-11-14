package line

import (
	"casino_backend/internal/repository"
	"errors"
	"sync"
)

type lineData struct {
	freeSpinsCount int
}

type repo struct {
	mtx   sync.RWMutex
	lines map[string]*lineData
}

func NewRepository() repository.LineRepository {
	return &repo{
		lines: make(map[string]*lineData),
	}
}

func (r *repo) GetCountFreeSpins(userID string) (int, error) {
	r.mtx.RLock()
	line, exists := r.lines[userID]
	r.mtx.RUnlock()

	if !exists {
		r.mtx.Lock()
		line, exists = r.lines[userID]
		if !exists {
			line = &lineData{freeSpinsCount: 0}
			r.lines[userID] = line
		}
		r.mtx.Unlock()
	}

	return line.freeSpinsCount, nil
}

func (r *repo) UpdateCountFreeSpins(userID string, newCount int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	line, exists := r.lines[userID]
	if !exists {
		return errors.New("line data not found")
	}

	line.freeSpinsCount = newCount
	return nil
}
