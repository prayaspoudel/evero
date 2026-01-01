package setup

import (
"fmt"

"github.com/spf13/viper"
)

// NewViper creates a new Viper configuration loader
// configPath: path to config directory (e.g., "config/healthcare")
// configName: name of config file without extension (e.g., "local", "development", "production")
func NewViper(configPath, configName string) *viper.Viper {
	config := viper.New()

	config.SetConfigName(configName)
	config.SetConfigType("json")
	config.AddConfigPath(configPath)
	config.AddConfigPath("./../../" + configPath)
	config.AddConfigPath("./../" + configPath)
	config.AddConfigPath("./" + configPath)
	
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w \n", err))
	}

	return config
}
