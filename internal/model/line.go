package model

type LineConfig struct {
	Reels    int
	Rows     int
	Lines    int
	Paylines [][]int

	SymbolWeights       map[string]int
	WildChanceOnReel234 float64
	Pays                map[string]map[int]int
	FreeSpinsByScatter  map[int]int
	BuyBonusMultiplier  int
	MaxPayoutMultiplier int
}

type SpinResult struct {
	Board            [5][3]rune
	LineWins         []LineWin
	ScatterCount     int
	ScatterPayout    int
	AwardedFreeSpins int
	TotalPayout      int
	Balance          int
	FreeSpinCount    int
	InFreeSpin       bool
}

type LineWin struct {
	Line   int
	Symbol rune
	Count  int
	Payout int
}
