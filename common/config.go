package common

import "github.com/spf13/viper"

var Cfg Config

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	JwtSecretKey  string `mapstructure:"JWT_SECRET_KEY"`
	IsProduction  bool   `mapstructure:"IS_PRODUCTION"`
	EncryptionKey string `mapstructure:"ENCRYPTION_KEY"`
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&Cfg)
	return
}
