package line

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
	err = s.repo.UpdateFreeSpinCount(10)
	return nil
}
