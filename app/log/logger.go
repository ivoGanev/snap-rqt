package logger

import (
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger  *log.Logger
	Println func(v ...any)
	Info    func(v ...any)
	Warning func(v ...any)
	Error   func(v ...any)
	Debug   func(v ...any)
)

func Init(logPath string) {
	// Truncate the log file before the logger starts writing
	err := os.Truncate(logPath, 0)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to truncate log file: %v", err)
	}

	logger = log.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   false,
	}, "", log.LstdFlags)

	Println = func(v ...any) {
		logger.Println(v...)
	}

	Info = func(v ...any) {
		logger.Println(append([]any{"[INFO]"}, sanitizeArgs(v)...)...)
	}

	Warning = func(v ...any) {
		logger.Println(append([]any{"[WARNING]"}, sanitizeArgs(v)...)...)
	}

	Error = func(v ...any) {
		logger.Println(append([]any{"[ERROR]"}, sanitizeArgs(v)...)...)
	}

	Debug = func(v ...any) {
		logger.Println(append([]any{"[DEBUG]"}, sanitizeArgs(v)...)...)
	}
}

func sanitizeArgs(args []any) []any {
	for i, arg := range args {
		if arg == nil {
			args[i] = "??? (nil)"
		} else if str, ok := arg.(string); ok && str == "" {
			args[i] = "??? (empty string)"
		}
	}
	return args
}
