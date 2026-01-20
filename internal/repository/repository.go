package repository

type LineRepository interface {
	GetBalance() (int, error)
	UpdateBalance(amount int) error
	GetFreeSpinCount() (int, error)
	UpdateFreeSpinCount(count int) error
}

type CascadeRepository interface {
	GetBalance() (int, error)
	UpdateBalance(amount int) error
	GetFreeSpinCount() (int, error)
	UpdateFreeSpinCount(count int) error

	GetMultiplierState() ([7][7]int, [7][7]int)
	SetMultiplierState(mult, hits [7][7]int) error
	ResetMultiplierState() error
}
