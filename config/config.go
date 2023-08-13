package config

import (
	"github.com/ory/viper"
)

type Config struct {
	HttpHost string `mapstructure:"HTTP_HOST"`
	HttpPort int    `mapstructure:"HTTP_PORT"`

	LogLevel int `mapstructure:"LOG_LEVEL"`

	MongoURI    string `mapstructure:"MONGO_URI"`
	MongoDBName string `mapstructure:"MONGO_DBNAME"`

	TokenAccessKey  string `mapstructure:"TOKEN_ACCESS_SECRET"`
	TokenRefreshKey string `mapstructure:"TOKEN_REFRESH_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
