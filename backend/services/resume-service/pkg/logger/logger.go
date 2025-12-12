package logger

import (
	"log"
	"os"
)

func NewLogger(serviceName, env string) *log.Logger {
	prefix := "[" + serviceName + "] "
	logger := log.New(os.Stdout, prefix, log.LstdFlags|log.Lshortfile)

	if env == "production" {
		logger.SetFlags(log.LstdFlags)
	}

	return logger
}
