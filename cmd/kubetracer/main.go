package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koustreak/kubetracer/internal/scanner"
	"github.com/koustreak/kubetracer/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	appName    = "kubetracer"
	appVersion = "1.0.0"
)

func main() {
	// Setup logger
	logger, err := utils.SetupLogger("info", "json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Starting KubeTracer")

	// Create Kubernetes client
	client, err := scanner.NewClient(logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create Kubernetes client")
	}

	// Create namespace scanner
	namespaceScanner := scanner.NewNamespaceScanner(client)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start background scanning routine
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		// Initial scan
		performNamespaceScan(ctx, namespaceScanner)

		// Periodic scans
		for {
			select {
			case <-ticker.C:
				performNamespaceScan(ctx, namespaceScanner)
			case <-ctx.Done():
				return
			}
		}
	}()

	logger.Info("KubeTracer is running")

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Shutting down")
	cancel()
}

// performNamespaceScan performs a namespace scan
func performNamespaceScan(ctx context.Context, namespaceScanner *scanner.NamespaceScanner) {
	namespaces, err := namespaceScanner.ListNamespaces(ctx)
	if err != nil {
		return
	}
	
	fmt.Printf("Found %d namespaces\n", len(namespaces))
}

// isSystemNamespace checks if a namespace is a system namespace
func isSystemNamespace(name string) bool {
	systemNamespaces := []string{
		"kube-system",
		"kube-public", 
		"kube-node-lease",
		"default",
	}
	
	for _, sysNs := range systemNamespaces {
		if name == sysNs {
			return true
		}
	}
	return false
}
