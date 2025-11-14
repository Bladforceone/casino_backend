package env

import (
	"casino_backend/internal/config"
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

const (
	symbolWeights = "SYMBOL_WEIGHTS"
	wildChance    = "WILD_CHANCE_ON_REEL_2_3_4"
)

type lineConfig struct {
	symbolWeights       map[string]int
	wildChanceOnReel234 float64
}

func NewLineConfig() (config.LineConfig, error) {
	symbolWeightsStr := os.Getenv(symbolWeights)
	if len(symbolWeightsStr) == 0 {
		return nil, errors.New("environment variable 'SYMBOL_WEIGHTS' is not set")
	}

	var symbolWeights map[string]int
	err := json.Unmarshal([]byte(symbolWeightsStr), &symbolWeights)
	if err != nil {
		return nil, err
	}

	wildChanceStr := os.Getenv(wildChance)
	if len(wildChanceStr) == 0 {
		return nil, errors.New("environment variable 'WILD_CHANCE' is not set")
	}

	wildChance, err := strconv.ParseFloat(wildChanceStr, 64)
	if err != nil {
		return nil, err
	}

	return &lineConfig{
		symbolWeights:       symbolWeights,
		wildChanceOnReel234: wildChance,
	}, nil
}

func (cfg *lineConfig) SymbolWeights() map[string]int {
	return cfg.symbolWeights
}

func (cfg *lineConfig) WildChance() float64 {
	return cfg.wildChanceOnReel234
}
