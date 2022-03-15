package loggerService

import (
	"github.com/mgutz/ansi"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
)

type LoggerService struct {
	// deps.
	io      *IOLogger
	manager *manager.Manager
	// vals.
	mode string
}

// NewLoggerService - constructor of LoggerService structure.
func NewLoggerService(manager *manager.Manager) *LoggerService {
	return &LoggerService{
		io:      NewIOLogger(),
		manager: manager,
		mode:    manager.Config.Environment.Mode,
	}
}

// Info - (!possibly panic!) log str into file with prefix INFO[dateTime]
func (logger *LoggerService) Info(str string) {
	logger.io.Write(str, logger.mode, LoggerInfoLevel)
}

// Debug - (!possibly panic!) log str into file with prefix DEBUG[dateTime]
func (logger *LoggerService) Debug(str string) {
	logger.io.Write(ansi.Color(str, "yellow"), logger.mode, LoggerDebugLevel)
}

// Error - (!possibly panic!) log str into file with prefix ERROR[dateTime]
func (logger *LoggerService) Error(str string) {
	logger.io.Write(ansi.Color(str, "red"), logger.mode, LoggerErrorLevel)
}

// Critical - (!possibly panic!) log str into file with prefix CRITICAL[dateTime]
func (logger *LoggerService) Critical(str string) {
	logger.io.Write(ansi.Color(str, "red+"), logger.mode, LoggerCriticalLevel)
}
