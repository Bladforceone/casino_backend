package repository

type UserRepository interface {
	GetBalance(userID string) (int, error)
	UpdateBalance(userID string, newBalance int) error
}

type LineRepository interface {
	GetCountFreeSpins(userID string) (int, error)
	UpdateCountFreeSpins(userID string, newCount int) error
}
