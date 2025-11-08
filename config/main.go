// Package config provides configuration management for the application using Viper.
package config

import "github.com/spf13/viper"

var config *viper.Viper
var AppConfig *viper.Viper

func init() {
	config = viper.New()

	config.SetConfigName("config")
	config.AddConfigPath("./")
	config.SetConfigType("yml")

	err := config.ReadInConfig()

	if err != nil {
		panic(err)
	}
}

func GetConfig() *viper.Viper {
	return config
}
