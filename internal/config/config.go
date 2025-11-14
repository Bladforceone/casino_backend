package config

type LineConfig interface {
	SymbolWeights() map[string]int
	WildChance() float64
}
