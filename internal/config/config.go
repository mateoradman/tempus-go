package config

import (
	"time"

	"github.com/spf13/viper"
)

// Configuration loaded from a .env file
type Config struct {
	DBSource             string        `mapstructure:"DB_URL"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	// SuperUser data
	SuperUserUsername string `mapstructure:"SUPERUSER_USERNAME"`
	SuperUserEmail    string `mapstructure:"SUPERUSER_EMAIL"`
	SuperUserPassword string `mapstructure:"SUPERUSER_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
