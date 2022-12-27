package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"reflect"
	"runtime"
	"strings"
)

const (
	envLogLevel  = "LOG_LEVEL"
	envLogOutput = "LOG_OUTPUT"
)

var (
	log logger
)

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type logger struct {
	log *zap.Logger
}

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutput()},
		Level:       zap.NewAtomicLevelAt(getLevel()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "time",
			MessageKey:   "msg",
			LevelKey:     "level",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log.log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

func getLevel() zapcore.Level {
	switch strings.ToLower(os.Getenv(envLogLevel)) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput() string {
	output := os.Getenv(envLogOutput)
	if output == "" {
		return "stdout"
	}
	return output
}

func Info(location interface{}, msg string, tags ...zap.Field) {
	switch location.(type) {
	case string:
		tags = append(tags, zap.String("location", location.(string)))
	default:
		tags = append(tags, zap.String("location", errorLocation(location)))
	}
	log.log.Info(msg, tags...)
	log.log.Sync()
}

func Error(location interface{}, err error, tags ...zap.Field) {
	switch location.(type) {
	case string:
		tags = append(tags, zap.String("location", location.(string)))
	default:
		tags = append(tags, zap.String("location", errorLocation(location)))
	}

	log.log.Error(err.Error(), tags...)
	log.log.Sync()
}

func Panic(location interface{}, err error, tags ...zap.Field) {
	switch location.(type) {
	case string:
		tags = append(tags, zap.String("location", location.(string)))
	default:
		tags = append(tags, zap.String("location", errorLocation(location)))
	}
	log.log.Panic(err.Error(), tags...)
	log.log.Sync()
}

func Warn(location interface{}, err error, tags ...zap.Field) {
	switch location.(type) {
	case string:
		tags = append(tags, zap.String("location", location.(string)))
	default:
		tags = append(tags, zap.String("location", errorLocation(location)))
	}
	log.log.Warn(err.Error(), tags...)
	log.log.Sync()
}

func errorLocation(temp interface{}) string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name(), ".")
	functionName := strs[len(strs)-1]
	strs = strings.Split(strs[len(strs)-2], "/")
	packageName := strs[len(strs)-1]
	return "ERROR_IN_" + strings.ToUpper(packageName) + "_" + strings.ToUpper(functionName)
}
