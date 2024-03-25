package core

import (
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	AlertProvider interface{}   `mapstructure:"alert_provider,omitempty"`
	DataProvider  interface{}   `mapstructure:"data_provider,omitempty"`
	Budgets       *[]Budget     `mapstructure:"budgets,omitempty"`
	Rules         *[]BudgetRule `mapstructure:"rules,omitempty"`
}

func LoadConfig(filePath string) (config Config, err error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return Config{}, err
	}
	log.Debug().
		Str("file", absPath).
		Msg("Reading config...")
	viper.SetConfigFile(filePath)
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&config)

	return
}

func SaveConfig(data interface{}) error {
	log.Debug().
		Str("path", viper.ConfigFileUsed()).
		Msg("Saving config...")
	return viper.WriteConfigAs(strings.Replace(viper.ConfigFileUsed(), ".wrangler", ".wranglerexp", 1))
}
