package initializers

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Database string `mapstructure:"DATABASE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		fmt.Println("Error reading config file", err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
