package pay

import (
	"casino_backend/internal/repository"
	"casino_backend/internal/service"
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

// Проверка соответствия интерфейсу
var _ service.PaymentService = (*serv)(nil)

type serv struct {
	txManager trm.Manager
	userRepo  repository.UserRepository
}

func NewService(
	txManager trm.Manager,
	userRepo repository.UserRepository,
) *serv {
	return &serv{
		txManager: txManager,
		userRepo:  userRepo,
	}
}

func (s *serv) Deposit(ctx context.Context, userID, amount int) error {
	return s.txManager.Do(ctx, func(txCtx context.Context) error {
		balance, err := s.userRepo.GetBalance(txCtx, userID)
		if err != nil {
			return err
		}

		newBalance := balance + amount
		return s.userRepo.UpdateBalance(txCtx, userID, newBalance)
	})
}

func (s *serv) GetBalance(ctx context.Context, userID int) (int, error) {
	return s.userRepo.GetBalance(ctx, userID)
}
