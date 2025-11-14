package env

import (
	"casino_backend/internal/config"
	"os"

	"gopkg.in/yaml.v3"
)

type lineConfig struct {
	SymbolWeightsData map[string]int `yaml:"symbol_weights"`
	WildChanceValue   float64        `yaml:"wild_chance_on_reel_2_3_4"`
	FreeSpinsScatter  map[int]int    `yaml:"free_spins_by_scatter"`
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
