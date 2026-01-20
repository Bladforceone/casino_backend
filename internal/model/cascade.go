package model

type CascadeSpin struct {
	Bet int
}

// Position представляет координаты ячейки на доске
type Position struct {
	Row int
	Col int
}

// ClusterInfo содержит информацию о найденном кластере
type ClusterInfo struct {
	Symbol     int        // Символ кластера
	Cells      []Position // Позиции ячеек в кластере
	Count      int        // Количество ячеек в кластере
	Payout     int        // Выигрыш за этот кластер в деньгах
	Multiplier int        // С каким итоговым множителем ушёл кластер
}

// CascadeStep представляет один шаг каскада
type CascadeStep struct {
	CascadeIndex int           // Номер каскада (0 - первый, 1 - второй и т.д.)
	Clusters     []ClusterInfo // Информация по всем кластерам на этом шаге
	NewSymbols   []struct {    // Это символы, которые упали сверху
		Position
		Symbol int
	}
}

// CascadeSpinResult представляет результат спина с каскадами
type CascadeSpinResult struct {
	InitialBoard     [7][7]int     // Доска сразу после начального заполнения, до любых каскадов
	Board            [7][7]int     // Итоговая доска после всех каскадов
	Cascades         []CascadeStep // Все шаги обновления доски
	TotalPayout      int           // Выигрыш за весь спин в деньгах
	Balance          int           // Баланс после спина в деньгах
	ScatterCount     int           // Количество бонусов, выпавших за спин
	AwardedFreeSpins int           // Количество начисленных фриспинов
	FreeSpinsLeft    int           // Остаток фриспинов после спина
	InFreeSpin       bool          // Находится ли игрок в режиме фриспинов
}

// CascadeData содержит информацию о балансе и количестве фриспинов игрока
type CascadeData struct {
	Balance       int // Теперь экспортировано (большая буква)
	FreeSpinCount int // Теперь экспортировано
}
