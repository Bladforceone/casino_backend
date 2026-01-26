package env

import (
	"casino_backend/internal/config"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type casdata struct {
	SymbolWeightsData map[int]int `yaml:"cascade_symbol_weights"`
	BonusPerColumn    float64     `yaml:"cascade_bonus_per_column"`
	BonusAwardsData   map[int]int `yaml:"cascade_bonus_awards"`
	PayTable          map[int]int `yaml:"cascade_pay_table"`
}

type cascadeConfigs struct {
	Configs []casdata `yaml:"configs"`
}

func NewCascadeConfigFromYAML(path string) (config.CascadeConfig, error) {
	confData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result cascadeConfigs
	if err := yaml.Unmarshal(confData, &result); err != nil {
		return nil, err
	}

	log.Println(result)
	return &result, nil
}

func (cfg *cascadeConfigs) SymbolWeights(idx int) map[int]int {
	return cfg.Configs[idx].SymbolWeightsData
}

func (cfg *cascadeConfigs) BonusProbPerColumn(idx int) float64 {
	return cfg.Configs[idx].BonusPerColumn
}

func (cfg *cascadeConfigs) BonusAwards(idx int) map[int]int {
	return cfg.Configs[idx].BonusAwardsData
}

func (cfg *cascadeConfigs) PayoutTable(idx int) map[int]int {
	return cfg.Configs[idx].PayTable
}
