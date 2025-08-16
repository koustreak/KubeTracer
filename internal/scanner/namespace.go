package scanner

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceScanner handles namespace scanning operations
type NamespaceScanner struct {
	client *Client
}

// NewNamespaceScanner creates a new namespace scanner
func NewNamespaceScanner(client *Client) *NamespaceScanner {
	return &NamespaceScanner{
		client: client,
	}
}

// ListNamespaces retrieves all namespaces from the cluster
func (ns *NamespaceScanner) ListNamespaces(ctx context.Context) ([]v1.Namespace, error) {
	ns.client.Logger.Info("Listing namespaces")
	
	namespaces, err := ns.client.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Log each namespace
	for _, namespace := range namespaces.Items {
		ns.client.Logger.WithFields(logrus.Fields{
			"namespace": namespace.Name,
			"status":    namespace.Status.Phase,
			"created":   namespace.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
		}).Info("Namespace found")
	}

	return namespaces.Items, nil
}
