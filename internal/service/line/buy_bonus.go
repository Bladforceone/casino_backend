package line

import (
	"casino_backend/internal/middleware"
	"context"
	"errors"
)

// BuyBonus Купить бонуску
func (s *serv) BuyBonus(ctx context.Context, amount int) error {
	cost := amount

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return errors.New("user id not found in context")
	}

	// Начало транзакции
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		balance, err := s.userRepo.GetBalance(txCtx, userID)
		if err != nil {
			return errors.New("failed to get user balance")
		}

		if balance < cost {
			return errors.New("not enough balance for bonus buy")
		}

		err = s.userRepo.UpdateBalance(txCtx, userID, balance-cost)
		if err != nil {
			return errors.New("failed to update user balance")
		}

		err = s.repo.UpdateFreeSpinCount(txCtx, userID, 10)
		if err != nil {
			return errors.New("failed to update free spin count")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
