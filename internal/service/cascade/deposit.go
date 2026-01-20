package cascade

// Пополнить баланс
func (s *serv) Deposit(amount int) error {
	err := s.repo.UpdateBalance(amount)
	if err != nil {
		return err
	}
	return nil
}
