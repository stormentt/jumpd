package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("jumpd")
	viper.AddConfigPath("/etc/jumpd")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No config file found. Using defaults")
	}

	viper.SetDefault("config.port", 22)
	viper.SetDefault("config.host", "127.0.0.1")
	viper.SetDefault("config.log.level", log.InfoLevel)

	viper.SetEnvPrefix("JUMPD")
	viper.AutomaticEnv()
}
