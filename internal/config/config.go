package config

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Scanner  ScannerConfig  `mapstructure:"scanner" yaml:"scanner"`
	Logging  LoggingConfig  `mapstructure:"logging" yaml:"logging"`
	Metrics  MetricsConfig  `mapstructure:"metrics" yaml:"metrics"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port         int    `mapstructure:"port" yaml:"port"`
	Host         string `mapstructure:"host" yaml:"host"`
	ReadTimeout  string `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout" yaml:"write_timeout"`
}

// ScannerConfig contains scanner-specific configuration
type ScannerConfig struct {
	Interval string `mapstructure:"interval" yaml:"interval"`
	Workers  int    `mapstructure:"workers" yaml:"workers"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level" yaml:"level"`
	Format string `mapstructure:"format" yaml:"format"`
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
	Port    int  `mapstructure:"port" yaml:"port"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Host:         "0.0.0.0",
			ReadTimeout:  "10s",
			WriteTimeout: "10s",
		},
		Scanner: ScannerConfig{
			Interval: "5m",
			Workers:  3,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Metrics: MetricsConfig{
			Enabled: true,
			Port:    9090,
		},
	}
}

// Load loads configuration from various sources
func Load() (*Config, error) {
	config := DefaultConfig()

	// Set configuration file paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/kubetracer/")
	viper.AddConfigPath("$HOME/.kubetracer")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("KUBETRACER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read configuration file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		logrus.Info("No config file found, using defaults and environment variables")
	}

	// Unmarshal configuration
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Scanner.Workers <= 0 {
		return fmt.Errorf("scanner workers must be greater than 0")
	}

	validLogLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	isValidLevel := false
	for _, level := range validLogLevels {
		if strings.ToLower(c.Logging.Level) == level {
			isValidLevel = true
			break
		}
	}
	if !isValidLevel {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	return nil
}
