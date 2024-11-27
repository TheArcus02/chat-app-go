package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return nil
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "", log.LstdFlags)
	
	return &Logger{logger}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log("FATAL", format, args...)
	os.Exit(1)
}

func (l *Logger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	message := fmt.Sprintf(format, args...)
	logMessage := fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)
	l.Logger.Println(logMessage)
}
