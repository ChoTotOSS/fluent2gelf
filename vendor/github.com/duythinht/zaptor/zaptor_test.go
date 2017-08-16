package zaptor

import (
	"os"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("testing")
	if logger.lvl.Level() != zapcore.InfoLevel {
		t.Log("Default log should be Info")
		t.Fail()
	}
}

func TestLogLevel(t *testing.T) {
	_ = os.Setenv("TEST_LEVEL", "DEBUG")
	logger := NewLogger("T").WithLevelBy("TEST_LEVEL")
	if logger.lvl.Level() != zapcore.DebugLevel {
		t.Log("Log level should be Debug")
		t.Fail()
	}
}
