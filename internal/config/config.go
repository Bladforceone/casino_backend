package config

type LineConfig interface {
	SymbolWeights() map[string]int
	WildChance() float64
	FreeSpinsByScatter() map[int]int
	PayoutTable() map[string]map[int]int
}

type CascadeConfig interface {
	SymbolWeights() map[int]int
	BonusProbPerColumn() float64
	BonusAwards() map[int]int
	PayoutTable() map[int]int
}
