package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Scanner handles Kubernetes cluster scanning operations
type Scanner struct {
	client kubernetes.Interface
	logger *logrus.Logger
}

// NewScanner creates a new scanner instance with in-cluster configuration
func NewScanner(logger *logrus.Logger) (*Scanner, error) {
	// Create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-cluster config: %w", err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Scanner{
		client: clientset,
		logger: logger,
	}, nil
}

// ListNamespaces retrieves all namespaces from the cluster
func (s *Scanner) ListNamespaces(ctx context.Context) ([]v1.Namespace, error) {
	s.logger.Info("Starting namespace discovery")
	
	startTime := time.Now()
	
	// List all namespaces
	namespaces, err := s.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		s.logger.WithError(err).Error("Failed to list namespaces")
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	duration := time.Since(startTime)
	
	s.logger.WithFields(logrus.Fields{
		"namespace_count": len(namespaces.Items),
		"duration_ms":     duration.Milliseconds(),
	}).Info("Successfully retrieved namespaces")

	// Log each namespace with details
	for _, ns := range namespaces.Items {
		s.logNamespaceDetails(ns)
	}

	return namespaces.Items, nil
}

// logNamespaceDetails logs detailed information about a namespace
func (s *Scanner) logNamespaceDetails(ns v1.Namespace) {
	fields := logrus.Fields{
		"namespace":     ns.Name,
		"status":        ns.Status.Phase,
		"created":       ns.CreationTimestamp.Format(time.RFC3339),
		"uid":          ns.UID,
	}

	// Add labels if present
	if len(ns.Labels) > 0 {
		fields["labels"] = ns.Labels
	}

	// Add annotations if present (excluding system annotations)
	if len(ns.Annotations) > 0 {
		userAnnotations := make(map[string]string)
		for key, value := range ns.Annotations {
			// Filter out system annotations to reduce noise
			if !isSystemAnnotation(key) {
				userAnnotations[key] = value
			}
		}
		if len(userAnnotations) > 0 {
			fields["annotations"] = userAnnotations
		}
	}

	s.logger.WithFields(fields).Info("Namespace discovered")
}

// isSystemAnnotation checks if an annotation is a system annotation
func isSystemAnnotation(key string) bool {
	systemPrefixes := []string{
		"kubectl.kubernetes.io/",
		"kubernetes.io/",
		"k8s.io/",
	}
	
	for _, prefix := range systemPrefixes {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// GetNamespaceCount returns the total number of namespaces
func (s *Scanner) GetNamespaceCount(ctx context.Context) (int, error) {
	namespaces, err := s.ListNamespaces(ctx)
	if err != nil {
		return 0, err
	}
	return len(namespaces), nil
}

// IsHealthy performs a health check by attempting to list namespaces
func (s *Scanner) IsHealthy(ctx context.Context) error {
	s.logger.Debug("Performing health check")
	
	// Simple health check - try to list namespaces
	_, err := s.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		s.logger.WithError(err).Error("Health check failed")
		return fmt.Errorf("health check failed: %w", err)
	}
	
	s.logger.Debug("Health check passed")
	return nil
}
