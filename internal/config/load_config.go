package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	DBParams struct {
		Host string `mapstructure:"host"`
		User string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name string `mapstructure:"name"`
		Port int `mapstructure:"port"`
	} `mapstructure:"db_params"`
	ServerParams struct {
		ServerAddress string `mapstructure:"server_address"`
		ListenPort string `mapstructure:"listen_port"`
		RateLimitMax int `mapstructure:"rate_limit_max"`
		RateLimitExpirationSeconds int `mapstructure:"rate_limit_expiration_seconds"`
	} `mapstructure:"server_params"`
}

func LoadConfig() *Config {
	viper.AddConfigPath("../internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	AppConfig := new(Config)
	viper.Unmarshal(AppConfig)
	return AppConfig
}
