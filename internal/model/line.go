package model

type LineSpin struct {
	Bet      int
	BuyBonus bool
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
