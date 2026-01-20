package cascade

import "casino_backend/internal/model"

func (s *serv) CheckData() (*model.CascadeData, error) {
	balance, err := s.repo.GetBalance()
	if err != nil {
		return nil, err
	}
	freeSpins, err := s.repo.GetFreeSpinCount()
	if err != nil {
		return nil, err
	}
	return &model.CascadeData{
		Balance:       balance, // Используем экспортированные имена
		FreeSpinCount: freeSpins,
	}, nil
}
