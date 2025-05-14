package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("\x1b[38;5;87m" + t.Format("02.01.2006 15:04:05.000") + "\x1b[0m")
	}

	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		color := ""
		text := l.CapitalString()

		switch l {
		case zapcore.DebugLevel:
			color = "\x1b[1;35m"
			text = "BYE"
		case zapcore.InfoLevel:
			color = "\x1b[1;34m"
		case zapcore.WarnLevel:
			color = "\x1b[1;33m"
		case zapcore.ErrorLevel:
			color = "\x1b[1;31m"
		case zapcore.FatalLevel:
			color = "\x1b[1;38;5;88m"
		}

		enc.AppendString(color + "[" + text + "]\x1b[0m")
	}

	config.DisableCaller = true
	config.DisableStacktrace = true

	logger, _ := config.Build()

	return logger
}

func main() {
	logger := initLogger()

	logger.Debug("Thanks for testing, goodbye!")
	logger.Info("hello")
	logger.Warn("hi?")
	logger.Error("oops..")
	//logger.Fatal("the end" + errors.New("fail").Error())

	fmt.Println(map[string]string{
		"hello": "ji",
		"hi":    "erwa",
	})
}

// 5 - 5.30: errors, logger, tokens(1)
// 5.30 - 6: tokens(2), tests
// 5 - 6: cobra, docker
// 6 - 7: readme, notice, docs
