package logger

import (
	"log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger  *log.Logger
	Println func(v ...any)
)

func Init(logPath string) {
	logger = log.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   false,
	}, "", log.LstdFlags)

	Println = logger.Println
	logger.Println("")
	logger.Println("---- Logger Init")
}
