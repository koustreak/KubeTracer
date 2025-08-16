package utils

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// SetupLogger configures and returns a structured logger
func SetupLogger(level, format string) (*logrus.Logger, error) {
	logger := logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		return nil, fmt.Errorf("invalid log level '%s': %w", level, err)
	}
	logger.SetLevel(logLevel)

	// Set formatter
	if strings.ToLower(format) == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return logger, nil
}
