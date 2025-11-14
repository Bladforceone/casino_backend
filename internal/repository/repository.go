package repository

type UserRepository interface {
	GetBalance(userID string) (int, error)
}

type LineRepository interface {
	GetCountFreeSpins(userID string) (int, error)
}
