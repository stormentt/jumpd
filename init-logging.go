package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initLogging() {
	logPath := viper.GetString("config.log.path")
	if len(logPath) > 0 {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("unable to open log file %s: %s\n", logPath, err)
		}

		log.SetOutput(f)
	}

	logLevel := viper.GetInt("config.log.level")
	switch logLevel {
	case 0:
		log.SetLevel(log.PanicLevel)
	case 1:
		log.SetLevel(log.FatalLevel)
	case 2:
		log.SetLevel(log.ErrorLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	case 4:
		log.SetLevel(log.InfoLevel)
	case 5:
		log.SetLevel(log.DebugLevel)
	}
}
