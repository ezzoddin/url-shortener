package database

import (
	"github.com/spf13/viper"
)

type Config struct {
	Driver  string `mapstructure:"DB_DRIVER"`
	Host    string `mapstructure:"DB_HOST"`
	User    string `mapstructure:"DB_USERNAME"`
	Pass    string `mapstructure:"DB_PASSWORD"`
	Port    string `mapstructure:"DB_PORT"`
	Name    string `mapstructure:"DB_DATABASE"`
	Sslmode string `mapstructure:"DB_SSLMODE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
