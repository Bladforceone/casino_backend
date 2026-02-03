package env

import (
	"casino_backend/internal/config"
	"os"

	"gopkg.in/yaml.v3"
)

type data struct {
	SymbolWeightsData map[string]int         `yaml:"line_symbol_weights"`
	WildChanceValue   float64                `yaml:"line_wild_chance_on_reel_2_3_4"`
	FreeSpinsScatter  map[int]int            `yaml:"line_free_spins_by_scatter"`
	PayTable          map[string]map[int]int `yaml:"line_payout_table"`
}

type lineConfig struct {
	Configs []data `yaml:"configs"`
}

func NewLineConfigFromYAML(path string) (config.LineConfig, error) {
	confData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result lineConfig
	if err := yaml.Unmarshal(confData, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (cfg *lineConfig) SymbolWeights(idx int) map[string]int {
	return cfg.Configs[idx].SymbolWeightsData
}

func (cfg *lineConfig) WildChance(idx int) float64 {
	return cfg.Configs[idx].WildChanceValue
}

func (cfg *lineConfig) FreeSpinsByScatter(idx int) map[int]int {
	return cfg.Configs[idx].FreeSpinsScatter
}

func (cfg *lineConfig) PayoutTable(idx int) map[string]map[int]int {
	return cfg.Configs[idx].PayTable
}
