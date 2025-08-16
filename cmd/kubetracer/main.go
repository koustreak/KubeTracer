package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koustreak/kubetracer/internal/config"
	"github.com/koustreak/kubetracer/internal/scanner"
	"github.com/koustreak/kubetracer/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	appName    = "kubetracer"
	appVersion = "1.0.0"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	logger, err := utils.SetupLogger(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	logger.WithFields(logrus.Fields{
		"app":     appName,
		"version": appVersion,
	}).Info("Starting KubeTracer")

	// Create scanner
	kubeScanner, err := scanner.NewScanner(logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create scanner")
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start background scanning routine
	go func() {
		interval, err := time.ParseDuration(cfg.Scanner.Interval)
		if err != nil {
			logger.WithError(err).Fatal("Invalid scanner interval")
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Initial scan
		if err := performNamespaceScan(ctx, kubeScanner, logger); err != nil {
			logger.WithError(err).Error("Initial namespace scan failed")
		}

		// Periodic scans
		for {
			select {
			case <-ticker.C:
				if err := performNamespaceScan(ctx, kubeScanner, logger); err != nil {
					logger.WithError(err).Error("Periodic namespace scan failed")
				}
			case <-ctx.Done():
				logger.Info("Stopping scanner due to context cancellation")
				return
			}
		}
	}()

	// Health check routine
	go func() {
		healthTicker := time.NewTicker(30 * time.Second)
		defer healthTicker.Stop()

		for {
			select {
			case <-healthTicker.C:
				if err := kubeScanner.IsHealthy(ctx); err != nil {
					logger.WithError(err).Warn("Health check failed")
				} else {
					logger.Debug("Health check passed")
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	logger.WithFields(logrus.Fields{
		"scan_interval": cfg.Scanner.Interval,
		"log_level":     cfg.Logging.Level,
	}).Info("KubeTracer is running")

	// Wait for shutdown signal
	sig := <-sigChan
	logger.WithField("signal", sig.String()).Info("Received shutdown signal")

	// Graceful shutdown
	logger.Info("Shutting down gracefully...")
	cancel()

	// Give some time for cleanup
	time.Sleep(2 * time.Second)
	logger.Info("Shutdown complete")
}

// performNamespaceScan performs a namespace scan and logs the results
func performNamespaceScan(ctx context.Context, kubeScanner *scanner.Scanner, logger *logrus.Logger) error {
	logger.Info("Starting namespace scan")
	
	startTime := time.Now()
	
	namespaces, err := kubeScanner.ListNamespaces(ctx)
	if err != nil {
		return fmt.Errorf("failed to scan namespaces: %w", err)
	}

	duration := time.Since(startTime)
	
	logger.WithFields(logrus.Fields{
		"total_namespaces": len(namespaces),
		"scan_duration":    duration.String(),
	}).Info("Namespace scan completed successfully")

	// Log summary
	systemNamespaces := 0
	userNamespaces := 0
	
	for _, ns := range namespaces {
		if isSystemNamespace(ns.Name) {
			systemNamespaces++
		} else {
			userNamespaces++
		}
	}

	logger.WithFields(logrus.Fields{
		"system_namespaces": systemNamespaces,
		"user_namespaces":   userNamespaces,
	}).Info("Namespace scan summary")

	return nil
}

// isSystemNamespace checks if a namespace is a system namespace
func isSystemNamespace(name string) bool {
	systemNamespaces := []string{
		"kube-system",
		"kube-public",
		"kube-node-lease",
		"local-path-storage",
		"default",
	}
	
	for _, sysNs := range systemNamespaces {
		if name == sysNs {
			return true
		}
	}
	
	// Check for system-like prefixes
	systemPrefixes := []string{
		"kube-",
		"kubernetes-",
	}
	
	for _, prefix := range systemPrefixes {
		if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
			return true
		}
	}
	
	return false
}
