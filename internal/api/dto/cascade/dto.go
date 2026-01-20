package cascade

type CascadeSpinRequest struct {
	Bet int `json:"bet"` // Размер ставки (положительное чётное число)
}

type CascadeSpinResponse struct {
	InitialBoard     [7][7]int     `json:"initial_board"`      // Доска сразу после начального заполнения, до любых каскадов
	Board            [7][7]int     `json:"board"`              // Итоговая доска: -1 = пусто, 0-6 = обычные, 7 = скаттер
	Cascades         []CascadeStep `json:"cascades"`           // Все шаги каскада (для анимации)
	TotalPayout      int           `json:"total_payout"`       // Общая выплата за спин
	Balance          int           `json:"balance"`            // Баланс после спина
	ScatterCount     int           `json:"scatter_count"`      // Количество скаттеров на финальной доске
	AwardedFreeSpins int           `json:"awarded_free_spins"` // Начислено фриспинов в этом спине
	FreeSpinsLeft    int           `json:"free_spins_left"`    // Остаток фриспинов после спина
	InFreeSpin       bool          `json:"in_free_spin"`       // Это был фриспин?
}

type CascadeStep struct {
	CascadeIndex int           `json:"cascade_index"` // 0 = первый, 1 = второй и т.д.
	Clusters     []ClusterInfo `json:"clusters"`      // Какие кластеры взорвались на этом шаге
	NewSymbols   []NewSymbol   `json:"new_symbols"`   // Новые символы, упавшие сверху
}

type ClusterInfo struct {
	Symbol     int        `json:"symbol"`     // ID символа (0–6)
	Cells      []Position `json:"cells"`      // Координаты ячеек в кластере
	Count      int        `json:"count"`      // Размер кластера (≥5)
	Payout     int        `json:"payout"`     // Выплата за кластер (в деньгах)
	Multiplier int        `json:"multiplier"` // Средний множитель (x2, x4, ..., x128)
}

type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type NewSymbol struct {
	Position Position `json:"position"`
	Symbol   int      `json:"symbol"` // -1 = пусто, 0–6 = обычный, 7 = скаттер
}

// Bonus Buy
type BuyCascadeBonusRequest struct {
	Amount int `json:"amount"` // Сумма покупки (обычно = bet × 100)
}

type BuyBonusResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message,omitempty"`
	AwardedSpins  int    `json:"awarded_spins,omitempty"`
	Cost          int    `json:"cost,omitempty"`
	Balance       int    `json:"balance,omitempty"`
	FreeSpinsLeft int    `json:"free_spins_left,omitempty"`
}

// Общий ответ на запрос данных (баланс + фриспины)
type CascadeDataResponse struct {
	Balance       int `json:"balance"`
	FreeSpinsLeft int `json:"free_spins_left"`
}
