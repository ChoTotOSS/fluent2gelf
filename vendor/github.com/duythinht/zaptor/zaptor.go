package zaptor

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	lvl zap.AtomicLevel
}

func NewLogger(name string) *Logger {

	level := zap.NewAtomicLevel()

	config := zap.Config{
		Level:            level,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()

	logger.Named(name)
	return &Logger{
		logger,
		level,
	}
}

func GetLogger(name string) *Logger {
	mx.Lock()
	defer mx.Unlock()
	if logger, ok := logs[name]; ok {
		return logger
	}

	logs[name] = NewLogger(name)
	return logs[name]
}

var (
	levels = map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info:": zap.InfoLevel,
		"error": zap.ErrorLevel,
		"fatal": zap.FatalLevel,
		"warn":  zap.WarnLevel,
	}

	logs = map[string]*Logger{}
	mx   = sync.Mutex{}
)

func (logger *Logger) WithLevelBy(envName string) *Logger {
	levelName := strings.ToLower(os.Getenv(envName))
	if lvl, ok := levels[levelName]; ok {
		logger.lvl.SetLevel(lvl)
	}
	return logger
}

func Default() *Logger {
	return NewLogger("default").WithLevelBy("LOG_LEVEL")
}
