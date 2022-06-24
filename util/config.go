package util

import (
	"github.com/spf13/viper"
)

// Config stores the application configuration from environment variables
type Config struct {
	DBUrl             string `mapstructure:"DB_URL"`
	DBDriver          string `mapstructure:"DB_DRIVER"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
}

// LoadConfig reads application configuration from environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
