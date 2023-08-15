package config

import (
	"log"

	"github.com/spf13/viper"
)

// Settings represents the application settings.
type Settings struct {
	FileName          string `mapstructure:"FILE_NAME"`
	ServerAddress     string `mapstructure:"SERVER_ADDRESS"`
	ServerPort        int    `mapstructure:"SERVER_PORT"`
	RedisAddress      string `mapstructure:"REDIS_ADDRESS"`
	RedisPort         int    `mapstructure:"REDIS_PORT"`
	RedisDB           int    `mapstructure:"REDIS_DB"`
	GinMode           string `mapstructure:"GIN_MODE"`
	ApiSecret         string `mapstructure:"API_SECRET"`
	ApiAuthSkipRoutes string `mapstructure:"API_AUTH_SKIP_ROUTES"`
}

// Load reads the configuration file and sets the application settings.
func Load() {
	viper.SetConfigName("quake")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Could not read config file: %s \n", err)
	}
}
