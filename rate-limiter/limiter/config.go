package limiter

import (
	"github.com/spf13/viper"
)

type conf struct {
	PeriodInSeconds string `mapstructure:"PERIOD"`
	LimitOfRequests string `mapstructure:"LIMIT"`
	BlockInSeconds  string `mapstructure:"BLOCK"`
	TokenLimit1     string `mapstructure:"LIMIT_TOKEN_HIGH"`
	TokenLimit2     string `mapstructure:"LIMIT_TOKEN_LOW"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
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
