package logger

import (
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger  *log.Logger
	Println func(v ...any)
)

func Init(logPath string) {
	// Truncate the log file before the logger starts writing
	err := os.Truncate(logPath, 0)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to truncate log file: %v", err)
	}

	logger = log.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   false,
	}, "", log.LstdFlags)

	Println = func(v ...any) {
		logger.Println(v...)
	}
}
