package logger

import (
	"io"
	"log"
	"os"
)

var (
	errorLogger *log.Logger
	warnLogger  *log.Logger
	infoLogger  *log.Logger
	debugLogger *log.Logger
	logLevel    = 1
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger = log.New(os.Stdout, "WARN ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(os.Stdout, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)
	if logLevel > 2 {
		warnLogger.SetOutput(io.Discard)
	}
	if logLevel > 1 {
		infoLogger.SetOutput(io.Discard)
	}
	if logLevel > 0 {
		debugLogger.SetOutput(io.Discard)
	}
}

func Info() *log.Logger {
	return infoLogger
}

func Debug() *log.Logger {
	return debugLogger
}

func Error() *log.Logger {
	return errorLogger
}

func Warn() *log.Logger {
	return warnLogger
}
