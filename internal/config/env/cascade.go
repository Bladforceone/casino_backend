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

type allCascadeConfigs struct {
	Configs []cascadeConfig `yaml:"configs"`
}

func NewCascadeConfigFromYAML(path string) ([]config.CascadeConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var ac allCascadeConfigs
	if err := yaml.Unmarshal(data, &ac); err != nil {
		return nil, err
	}
	var configs []config.CascadeConfig
	for i := range ac.Configs {
		configs = append(configs, &ac.Configs[i])
	}
	return configs, nil
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
