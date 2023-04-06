package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ServerAddress    string `mapstructure:"SERVER_ADDRESS"`
	RestaurantsTable string `mapstructure:"RESTAURANTS_TABLE"`
	PlaceIndex       string `mapstructure:"PLACE_INDEX"`
}

// Init reads configuration from file or environment variables.
func Init(path string) Config {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error on reading configuration file: %s", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("error on parsing configuration file: %s", err)
	}

	return config
}
