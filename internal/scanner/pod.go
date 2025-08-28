package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodInfo represents pod information structure
type PodInfo struct {
	Status       string `json:"status"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	CapturedTime string `json:"captured_time"`
}

// PodScanner handles pod scanning operations
type PodScanner struct {
	client *Client
}

// NewPodScanner creates a new pod scanner
func NewPodScanner(client *Client) *PodScanner {
	return &PodScanner{
		client: client,
	}
}

// ListPods retrieves all pods from the specified namespace
func (ps *PodScanner) ListPods(ctx context.Context, namespace string) (map[string]PodInfo, error) {
	ps.client.Logger.WithFields(logrus.Fields{
		"namespace": namespace,
	}).Info("Listing pods in namespace")

	pods, err := ps.client.Clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	// Create result map
	result := make(map[string]PodInfo)
	capturedTime := time.Now().Format("02-Jan-2006 15:04:05")

	// Process each pod
	for _, pod := range pods.Items {
		podInfo := PodInfo{
			Status:       string(pod.Status.Phase),
			StartTime:    "",
			EndTime:      "",
			CapturedTime: capturedTime,
		}

		// Set start time if available
		if pod.Status.StartTime != nil {
			podInfo.StartTime = pod.Status.StartTime.Format("02-Jan-2006 15:04:05")
		}

		// Set end time if pod has finished (Failed, Succeeded)
		if pod.Status.Phase == v1.PodFailed || pod.Status.Phase == v1.PodSucceeded {
			// Look for container statuses to get end time
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.State.Terminated != nil {
					podInfo.EndTime = containerStatus.State.Terminated.FinishedAt.Format("02-Jan-2006 15:04:05")
					break
				}
			}
		}

		// Add to result map
		result[pod.Name] = podInfo

		// Log pod information
		ps.client.Logger.WithFields(logrus.Fields{
			"pod":       pod.Name,
			"namespace": namespace,
			"status":    podInfo.Status,
			"startTime": podInfo.StartTime,
			"endTime":   podInfo.EndTime,
			"captured":  podInfo.CapturedTime,
		}).Info("Pod found")
	}

	ps.client.Logger.WithFields(logrus.Fields{
		"namespace": namespace,
		"podCount":  len(result),
	}).Info("Pod listing completed")

	return result, nil
}

// ListAllPods retrieves pods from all namespaces (if no specific namespace provided)
func (ps *PodScanner) ListAllPods(ctx context.Context) (map[string]map[string]PodInfo, error) {
	ps.client.Logger.Info("Listing pods from all namespaces")

	// First get all namespaces
	namespaces, err := ps.client.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Result map: namespace -> pod name -> pod info
	result := make(map[string]map[string]PodInfo)

	// Get pods from each namespace
	for _, ns := range namespaces.Items {
		pods, err := ps.ListPods(ctx, ns.Name)
		if err != nil {
			ps.client.Logger.WithFields(logrus.Fields{
				"namespace": ns.Name,
				"error":     err,
			}).Warn("Failed to list pods in namespace")
			continue
		}

		if len(pods) > 0 {
			result[ns.Name] = pods
		}
	}

	ps.client.Logger.WithFields(logrus.Fields{
		"namespaceCount": len(result),
	}).Info("All pods listing completed")

	return result, nil
}
