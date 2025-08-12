package db

import (
	"fmt"
	"github.com/lzkking/harle/server/config"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

type gormWriter struct {
}

func (w gormWriter) Printf(format string, args ...interface{}) {
	zap.S().Infof(fmt.Sprintf(format, args...))
}

func getGormLogger(dbConfig *config.DatabaseConfig) logger.Interface {
	logConfig := logger.Config{
		SlowThreshold: time.Second,
		Colorful:      true,
		LogLevel:      logger.Info,
	}
	switch strings.ToLower(dbConfig.LogLevel) {
	case "silent":
		logConfig.LogLevel = logger.Silent
	case "err":
		fallthrough
	case "error":
		logConfig.LogLevel = logger.Error
	case "warning":
		fallthrough
	case "warn":
		logConfig.LogLevel = logger.Warn
	case "info":
		fallthrough
	default:
		logConfig.LogLevel = logger.Info
	}

	return logger.New(gormWriter{}, logConfig)
}
