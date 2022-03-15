package loggerService

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
)

const LoggerInfoLevel = "info"
const LoggerDebugLevel = "debug"
const LoggerErrorLevel = "error"
const LoggerCriticalLevel = "critical"

var availableModesPathsMap = map[string]string{
	config.DevMode:  "/app/var/log/dev.log",
	config.ProdMode: "/app/var/log/prod.log",
}

var availableLoggerLevelsMap = map[string]string{
	LoggerInfoLevel:     LoggerInfoLevel,
	LoggerDebugLevel:    LoggerDebugLevel,
	LoggerErrorLevel:    LoggerErrorLevel,
	LoggerCriticalLevel: LoggerCriticalLevel,
}

type IOLoggerInterface interface {
	// The only one method, which must be implemented by other instance of logger.
	Write(logStr string, mode string, level string)
}

// Finall structure, can be instantinate properly from constructor in any place of code.
type IOLogger struct {
}

// NewIOLogger - constructor of IOLogger structure.
func NewIOLogger() *IOLogger {
	return &IOLogger{}
}

// Write - write string into log file.
func (io *IOLogger) Write(logStr string, mode string, level string) {
	path, ok := availableModesPathsMap[mode]
	level, lvlOk := availableLoggerLevelsMap[level]

	var f *os.File
	var err error
	if ok && lvlOk {
		if mode == config.DevMode {
			f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		} else {
			f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		}
		if err != nil {
			panic(fmt.Sprintf("IOLogger: error opening file: %v!", err))
		}
		defer f.Close()

		log.Println(logStr)
		log.SetOutput(f)

		// writing log to the destination
		wasWritten := io.write(logStr, level)
		if !wasWritten {
			panic("The log is not written. The levels didn't match.")
		}

		return
	}

	panic("IOLogger: received unexpected environment mode!")
}

// write - write the logStr to destination log file path.
func (io *IOLogger) write(logStr string, level string) bool {
	lvl, ok := availableLoggerLevelsMap[level]
	if ok {
		log.Print(
			fmt.Sprintf(
				"%v: %v\n\n",
				strings.ToUpper(lvl),
				logStr,
			),
		)

		return true
	}

	return false
}

// FileExists - check the target file path is exists.
func (io *IOLogger) FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// CreateFile - create a file or truncate if it is exists with 0666 mask.
func (io *IOLogger) CreateFile(path string) error {
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		fo.Close()
	}()

	return nil
}
