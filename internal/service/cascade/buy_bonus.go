package cascade

import (
	"errors"
)

// Купить бонуску
func (s *serv) BuyBonus(amount int) error {
	cost := amount

	balance, err := s.repo.GetBalance()
	if err != nil {
		return errors.New("failed to get user balance")
	}
	if balance < cost {
		return errors.New("not enough balance for bonus buy")
	}
	err = s.repo.UpdateBalance(balance - cost)
	if err != nil {
		return errors.New("failed to update balance after bonus buy")
	}

	if err := s.repo.ResetMultiplierState(); err != nil {
		return errors.New("failed to reset mult state")
	}

	err = s.repo.UpdateFreeSpinCount(10)
	if err != nil {
		return errors.New("failed to update free spin count after bonus buy")
	}
	return nil
}
