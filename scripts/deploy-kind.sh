#!/bin/bash

# Build and deployment script for KubeTracer in KIND cluster

set -e

# Configuration
APP_NAME="kubetracer"
DOCKER_USER="kdutta93"
DOCKER_IMAGE="$DOCKER_USER/kubetracer"
HELM_CHART_PATH="./deployments/helm/kubetracer"
NAMESPACE="kubetracer"

echo "🚀 Starting KubeTracer deployment process..."

# Step 1: Create namespace
echo "🏗️  Creating namespace..."
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# Step 2: Deploy using Helm (pulls image from Docker Hub)
echo "🎯 Deploying KubeTracer using Helm..."
helm upgrade --install "$APP_NAME" "$HELM_CHART_PATH" \
    --namespace "$NAMESPACE" \
    --set image.repository="$DOCKER_IMAGE" \
    --set image.tag="latest" \
    --set image.pullPolicy="Always" \
    --wait --timeout=300s

# Step 3: Check deployment status
echo "✅ Checking deployment status..."
kubectl get pods -n "$NAMESPACE"
kubectl get services -n "$NAMESPACE"

echo "🎉 Deployment completed successfully!"
echo ""
echo "To check logs, run:"
echo "  kubectl logs -f deployment/$APP_NAME -n $NAMESPACE"
echo ""
echo "To port-forward the service, run:"
echo "  kubectl port-forward service/$APP_NAME 8080:8080 -n $NAMESPACE"
