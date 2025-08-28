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

	// Create pod scanner
	podScanner := scanner.NewPodScanner(client)

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
		performPodScan(ctx, podScanner, "kube-system") // Example: scan kube-system namespace
		// performAllPodsScan(ctx, podScanner) // Uncomment to scan all namespaces

		// Periodic scans
		for {
			select {
			case <-ticker.C:
				performNamespaceScan(ctx, namespaceScanner)
				performPodScan(ctx, podScanner, "kube-system") // Example: scan kube-system namespace
				// performAllPodsScan(ctx, podScanner) // Uncomment to scan all namespaces
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

// performPodScan performs a pod scan for a specific namespace
// Usage: performPodScan(ctx, podScanner, "kube-system")
// Returns: map[string]PodInfo where key is pod name and value contains:
//   - Status: Pod phase (Running, Pending, Succeeded, Failed, Unknown)
//   - StartTime: Pod start time in "DD-Mon-YYYY HH:MM:SS" format
//   - EndTime: Pod end time (empty if still running)
//   - CapturedTime: Current time when data was captured
func performPodScan(ctx context.Context, podScanner *scanner.PodScanner, namespace string) {
	pods, err := podScanner.ListPods(ctx, namespace)
	if err != nil {
		fmt.Printf("Error scanning pods in namespace %s: %v\n", namespace, err)
		return
	}

	fmt.Printf("Found %d pods in namespace %s\n", len(pods), namespace)

	// Print pod details
	for podName, podInfo := range pods {
		fmt.Printf("Pod: %s - Status: %s, Start: %s, End: %s, Captured: %s\n",
			podName, podInfo.Status, podInfo.StartTime, podInfo.EndTime, podInfo.CapturedTime)
	}
}

// performAllPodsScan performs a pod scan across all namespaces
func performAllPodsScan(ctx context.Context, podScanner *scanner.PodScanner) {
	allPods, err := podScanner.ListAllPods(ctx)
	if err != nil {
		fmt.Printf("Error scanning all pods: %v\n", err)
		return
	}

	totalPods := 0
	for namespace, pods := range allPods {
		fmt.Printf("Namespace: %s - %d pods\n", namespace, len(pods))
		totalPods += len(pods)

		// Print pod details for each namespace
		for podName, podInfo := range pods {
			fmt.Printf("  Pod: %s - Status: %s, Start: %s, End: %s\n",
				podName, podInfo.Status, podInfo.StartTime, podInfo.EndTime)
		}
	}

	fmt.Printf("Total pods found across all namespaces: %d\n", totalPods)
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
