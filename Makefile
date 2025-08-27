# Build configuration
APP_NAME = kubetracer
DOCKER_IMAGE = kdutta93/kubetracer:latest
DOCKER_USER = kdutta93
HELM_CHART = deployments/helm/kubetracer
NAMESPACE = kubetracer

# Go configuration
GO_VERSION = 1.22
GOOS = linux
GOARCH = amd64

# Build the Go binary
.PHONY: build
build:
	@echo "🔨 Building $(APP_NAME)..."
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo -o bin/$(APP_NAME) ./cmd/kubetracer

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	docker rmi $(DOCKER_IMAGE) 2>/dev/null || true

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_USER)/$(APP_NAME):latest .

# Push Docker image to Docker Hub
.PHONY: docker-push
docker-push: docker-build
	@echo "📤 Pushing Docker image to Docker Hub..."
	docker push $(DOCKER_USER)/$(APP_NAME):latest

# Deploy to KIND cluster (pulls from Docker Hub)
.PHONY: deploy-kind
deploy-kind: docker-push
	@echo "🚀 Deploying to KIND cluster..."
	./scripts/deploy-kind.sh

# Deploy only (assumes image is already pushed to Docker Hub)
.PHONY: deploy-only
deploy-only:
	@echo "🚀 Deploying to KIND cluster (skipping build/push)..."
	./scripts/deploy-kind.sh

# Install dependencies
.PHONY: deps
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the Go binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-push  - Build and push Docker image to Docker Hub"
	@echo "  deploy-kind  - Build, push and deploy to KIND cluster"
	@echo "  deploy-only  - Deploy to KIND cluster (skip build/push)"
	@echo "  deps         - Install dependencies"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Usage examples:"
	@echo "  make deploy-kind                    # Use default 'kubetracer' Docker Hub user"
	@echo "  DOCKER_USER=myuser make deploy-kind # Use custom Docker Hub username"

.DEFAULT_GOAL := help
