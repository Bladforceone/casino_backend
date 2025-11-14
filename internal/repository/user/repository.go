package user

import (
	"casino_backend/internal/repository"
	"errors"
	"sync"
)

type userData struct {
	id       string
	email    string
	password string
	balance  int
}

type repo struct {
	mtx   sync.RWMutex
	users map[string]*userData
}

func NewRepository() repository.UserRepository {
	return &repo{
		users: make(map[string]*userData),
	}
}

func (repo *repo) GetBalance(userID string) (int, error) {
	repo.mtx.RLock()
	defer repo.mtx.RUnlock()
	user, exists := repo.users[userID]
	if !exists {
		return 0, errors.New("user not found")
	}

	return user.balance, nil
}
