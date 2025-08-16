package scanner

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client wraps the Kubernetes client with logging
type Client struct {
	Clientset kubernetes.Interface
	Logger    *logrus.Logger
}

// NewClient creates a new Kubernetes client with in-cluster configuration
func NewClient(logger *logrus.Logger) (*Client, error) {
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

	return &Client{
		Clientset: clientset,
		Logger:    logger,
	}, nil
}
