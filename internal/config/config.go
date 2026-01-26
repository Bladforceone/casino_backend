package config

import (
	"time"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	return nil
}

type LineConfig interface {
	SymbolWeights(idx int) map[string]int
	WildChance(idx int) float64
	FreeSpinsByScatter(idx int) map[int]int
	PayoutTable(idx int) map[string]map[int]int
}

type CascadeConfig interface {
	SymbolWeights(idx int) map[int]int
	BonusProbPerColumn(idx int) float64
	BonusAwards(idx int) map[int]int
	PayoutTable(idx int) map[int]int
}

type HTTPConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type JWTConfig interface {
	AccessTokenSecretKey() []byte
	AccessTokenDuration() time.Duration
	RefreshTokenDuration() time.Duration
}
