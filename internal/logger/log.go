package logger

import (
	"go.uber.org/zap"
	"log"
	"sync"
)

var loggerInstance *zap.Logger
var once sync.Once

func InitLogger() {
	once.Do(func() {
		GetLogger()
	})
}

func GetLogger() *zap.Logger {
	if loggerInstance == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			log.Fatal("Unable to create logger")
			return nil
		}
		loggerInstance = logger
	}
	return loggerInstance
}

func Close() {
	if loggerInstance != nil {
		loggerInstance.Sync()
		loggerInstance = nil
	}
}
