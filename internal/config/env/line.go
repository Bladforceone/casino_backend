package env

import (
	"casino_backend/internal/config"
	"os"

	"gopkg.in/yaml.v3"
)

type lineConfig struct {
	SymbolWeightsData map[string]int         `yaml:"line_symbol_weights"`
	WildChanceValue   float64                `yaml:"line_wild_chance_on_reel_2_3_4"`
	FreeSpinsScatter  map[int]int            `yaml:"line_free_spins_by_scatter"`
	PayTable          map[string]map[int]int `yaml:"line_payout_table"`
}

func NewLineConfigFromYAML(path string) (config.LineConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg lineConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *lineConfig) SymbolWeights() map[string]int {
	return cfg.SymbolWeightsData
}

func (cfg *lineConfig) WildChance() float64 {
	return cfg.WildChanceValue
}

func (cfg *lineConfig) FreeSpinsByScatter() map[int]int {
	return cfg.FreeSpinsScatter
}

func (cfg *lineConfig) PayoutTable() map[string]map[int]int {
	return cfg.PayTable
}
