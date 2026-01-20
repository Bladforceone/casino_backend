package env

import (
	"casino_backend/internal/config"
	"os"

	"gopkg.in/yaml.v3"
)

type cascadeConfig struct {
	SymbolWeightsData map[int]int `yaml:"cascade_symbol_weights"`
	BonusPerColumn    float64     `yaml:"cascade_bonus_per_column"`
	BonusAwardsData   map[int]int `yaml:"cascade_bonus_awards"`
	PayTable          map[int]int `yaml:"cascade_pay_table"`
}

func NewCascadeConfigFromYAML(path string) (config.CascadeConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg cascadeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *cascadeConfig) SymbolWeights() map[int]int {
	return cfg.SymbolWeightsData
}

func (cfg *cascadeConfig) BonusProbPerColumn() float64 {
	return cfg.BonusPerColumn
}

func (cfg *cascadeConfig) BonusAwards() map[int]int {
	return cfg.BonusAwardsData
}

func (cfg *cascadeConfig) PayoutTable() map[int]int {
	return cfg.PayTable
}
