#!/bin/bash

# Build and deployment script for KubeTracer in KIND cluster

set -e

# Configuration
APP_NAME="kubetracer"
DOCKER_IMAGE="kubetracer/kubetracer:latest"
HELM_CHART_PATH="./deployments/helm/kubetracer"
NAMESPACE="kubetracer"

echo "🚀 Starting KubeTracer build and deployment process..."

# Step 1: Build Docker image
echo "📦 Building Docker image..."
docker build -t "$DOCKER_IMAGE" .

# Step 2: Load image into KIND cluster
echo "📋 Loading image into KIND cluster..."
kind load docker-image "$DOCKER_IMAGE"

# Step 3: Create namespace
echo "🏗️  Creating namespace..."
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# Step 4: Deploy using Helm
echo "🎯 Deploying KubeTracer using Helm..."
helm upgrade --install "$APP_NAME" "$HELM_CHART_PATH" \
    --namespace "$NAMESPACE" \
    --set image.repository="kubetracer/kubetracer" \
    --set image.tag="latest" \
    --set image.pullPolicy="Never" \
    --wait --timeout=300s

# Step 5: Check deployment status
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
