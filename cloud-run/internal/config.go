package internal

import (
	"github.com/spf13/viper"
)

type conf struct {
	Key string `mapstructure:"KEY"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	//viper.AddConfigPath("/home/susin/code/gsusin/goexpert/rate-limiter/cmd")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
