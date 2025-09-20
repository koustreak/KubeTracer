package scanner

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretInfo represents secret information structure
type SecretInfo struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	DataKeys []string          `json:"data_keys"`
	Created  string            `json:"created"`
	Labels   map[string]string `json:"labels,omitempty"`
}

// SecretScanner handles secret scanning operations
type SecretScanner struct {
	client *Client
}

// NewSecretScanner creates a new secret scanner
func NewSecretScanner(client *Client) *SecretScanner {
	return &SecretScanner{
		client: client,
	}
}

// ListSecrets retrieves all secrets from the specified namespace
func (ss *SecretScanner) ListSecrets(ctx context.Context, namespace string) ([]SecretInfo, error) {
	ss.client.Logger.WithFields(logrus.Fields{
		"namespace": namespace,
	}).Info("Listing secrets in namespace")

	secrets, err := ss.client.Clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets in namespace %s: %w", namespace, err)
	}

	// Create result slice
	var result []SecretInfo

	// Process each secret
	for _, secret := range secrets.Items {
		// Get data keys (without exposing actual secret values)
		dataKeys := make([]string, 0, len(secret.Data))
		for key := range secret.Data {
			dataKeys = append(dataKeys, key)
		}

		secretInfo := SecretInfo{
			Name:     secret.Name,
			Type:     string(secret.Type),
			DataKeys: dataKeys,
			Created:  secret.CreationTimestamp.Format("02-Jan-2006 15:04:05"),
			Labels:   secret.Labels,
		}

		result = append(result, secretInfo)

		// Log secret information (without sensitive data)
		ss.client.Logger.WithFields(logrus.Fields{
			"secret":    secret.Name,
			"namespace": namespace,
			"type":      secretInfo.Type,
			"dataKeys":  len(dataKeys),
			"created":   secretInfo.Created,
		}).Info("Secret found")
	}

	ss.client.Logger.WithFields(logrus.Fields{
		"namespace":   namespace,
		"secretCount": len(result),
	}).Info("Secret listing completed")

	return result, nil
}

// ListAllSecrets retrieves secrets from all namespaces
func (ss *SecretScanner) ListAllSecrets(ctx context.Context) (map[string][]SecretInfo, error) {
	ss.client.Logger.Info("Listing secrets from all namespaces")

	// First get all namespaces
	namespaces, err := ss.client.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Result map: namespace -> list of secrets
	result := make(map[string][]SecretInfo)

	// Get secrets from each namespace
	for _, ns := range namespaces.Items {
		secrets, err := ss.ListSecrets(ctx, ns.Name)
		if err != nil {
			ss.client.Logger.WithFields(logrus.Fields{
				"namespace": ns.Name,
				"error":     err,
			}).Warn("Failed to list secrets in namespace")
			continue
		}

		// Only add namespace to result if it has secrets
		if len(secrets) > 0 {
			result[ns.Name] = secrets
		}
	}

	totalSecrets := 0
	for _, secrets := range result {
		totalSecrets += len(secrets)
	}

	ss.client.Logger.WithFields(logrus.Fields{
		"namespaceCount": len(result),
		"totalSecrets":   totalSecrets,
	}).Info("All secrets listing completed")

	return result, nil
}

// ListSecretsByType retrieves secrets of a specific type from all namespaces
func (ss *SecretScanner) ListSecretsByType(ctx context.Context, secretType string) (map[string][]SecretInfo, error) {
	ss.client.Logger.WithFields(logrus.Fields{
		"secretType": secretType,
	}).Info("Listing secrets by type from all namespaces")

	// Get all secrets first
	allSecrets, err := ss.ListAllSecrets(ctx)
	if err != nil {
		return nil, err
	}

	// Filter by secret type
	result := make(map[string][]SecretInfo)

	for namespace, secrets := range allSecrets {
		var filteredSecrets []SecretInfo
		for _, secret := range secrets {
			if secret.Type == secretType {
				filteredSecrets = append(filteredSecrets, secret)
			}
		}
		if len(filteredSecrets) > 0 {
			result[namespace] = filteredSecrets
		}
	}

	totalFiltered := 0
	for _, secrets := range result {
		totalFiltered += len(secrets)
	}

	ss.client.Logger.WithFields(logrus.Fields{
		"secretType":      secretType,
		"namespaceCount":  len(result),
		"filteredSecrets": totalFiltered,
	}).Info("Secrets filtering by type completed")

	return result, nil
}
