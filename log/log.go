package log

import (
	"github.com/charmbracelet/log"
)

var logger *MyLogger = nil

func GetLogger() *MyLogger {
	if logger != nil {
		return logger
	}

	logger = log.Default()
	logger.SetTimeFormat("15:04:05.000")

	return logger
}

type MyLogger = log.Logger
