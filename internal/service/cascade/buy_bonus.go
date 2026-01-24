package cascade

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

	balance, err := s.userRepo.GetBalance(ctx, userID)
	if err != nil {
		return errors.New("failed to get user balance")
	}
	if balance < cost {
		return errors.New("not enough balance for bonus buy")
	}
	err = s.userRepo.UpdateBalance(ctx, userID, balance-cost)
	if err != nil {
		return errors.New("failed to update balance after bonus buy")
	}

	if err := s.cascadeRepo.ResetMultiplierState(ctx, userID); err != nil {
		return errors.New("failed to reset mult state")
	}

	err = s.cascadeRepo.UpdateFreeSpinCount(ctx, userID, 10)
	if err != nil {
		return errors.New("failed to update free spin count after bonus buy")
	}
	return nil
}
