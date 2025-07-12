# KubeTracer Makefile
# Variables
APP_NAME := kubetracer
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT := $(shell git rev-parse HEAD)
GO_VERSION := $(shell go version | awk '{print $$3}')

# Build variables
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.gitCommit=$(GIT_COMMIT)"
BINARY_DIR := bin
DOCKER_REGISTRY := your-registry.com
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(APP_NAME)

# Go related variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt

# Default target
.PHONY: all
all: clean test build

# Help target
.PHONY: help
help: ## Display this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
.PHONY: build
build: build-cli build-operator build-exporter ## Build all binaries

.PHONY: build-cli
build-cli: ## Build CLI binary
	@echo "Building CLI..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(APP_NAME) ./cmd/cli

.PHONY: build-operator
build-operator: ## Build operator binary
	@echo "Building operator..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(APP_NAME)-operator ./cmd/operator

.PHONY: build-exporter
build-exporter: ## Build exporter binary
	@echo "Building exporter..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(APP_NAME)-exporter ./cmd/exporter

.PHONY: build-linux
build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(APP_NAME)-linux-amd64 ./cmd/cli

.PHONY: build-windows
build-windows: ## Build for Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/cli

.PHONY: build-darwin
build-darwin: ## Build for macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/cli

.PHONY: build-all-platforms
build-all-platforms: build-linux build-windows build-darwin ## Build for all platforms

# Development targets
.PHONY: run
run: ## Run the CLI application
	$(GOCMD) run ./cmd/cli

.PHONY: run-operator
run-operator: ## Run the operator
	$(GOCMD) run ./cmd/operator

.PHONY: install
install: ## Install the CLI binary
	$(GOCMD) install $(LDFLAGS) ./cmd/cli

# Testing targets
.PHONY: test
test: ## Run unit tests
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

.PHONY: test-integration
test-integration: ## Run integration tests
	$(GOTEST) -v -tags=integration ./test/integration/...

.PHONY: test-e2e
test-e2e: ## Run end-to-end tests
	$(GOTEST) -v -tags=e2e ./test/e2e/...

.PHONY: benchmark
benchmark: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

# Code quality targets
.PHONY: fmt
fmt: ## Format Go code
	$(GOFMT) -s -w .

.PHONY: vet
vet: ## Run go vet
	$(GOCMD) vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

.PHONY: security
security: ## Run security checks with gosec
	gosec ./...

.PHONY: check
check: fmt vet lint security ## Run all code quality checks

# Dependency management
.PHONY: deps
deps: ## Download dependencies
	$(GOMOD) download

.PHONY: deps-update
deps-update: ## Update dependencies
	$(GOMOD) tidy
	$(GOGET) -u ./...

.PHONY: deps-vendor
deps-vendor: ## Vendor dependencies
	$(GOMOD) vendor

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(VERSION) -f deployments/docker/Dockerfile .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

.PHONY: docker-push
docker-push: ## Push Docker image
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run --rm -it $(DOCKER_IMAGE):latest

# Kubernetes targets
.PHONY: k8s-deploy
k8s-deploy: ## Deploy to Kubernetes
	kubectl apply -f deployments/k8s/

.PHONY: k8s-delete
k8s-delete: ## Delete from Kubernetes
	kubectl delete -f deployments/k8s/

.PHONY: helm-install
helm-install: ## Install Helm chart
	helm install $(APP_NAME) deployments/helm/$(APP_NAME)

.PHONY: helm-upgrade
helm-upgrade: ## Upgrade Helm chart
	helm upgrade $(APP_NAME) deployments/helm/$(APP_NAME)

.PHONY: helm-uninstall
helm-uninstall: ## Uninstall Helm chart
	helm uninstall $(APP_NAME)

# Protocol Buffers
.PHONY: proto
proto: ## Generate protobuf files
	protoc --go_out=. --go-grpc_out=. api/proto/*.proto

# Documentation
.PHONY: docs
docs: ## Generate documentation
	godoc -http=:6060

.PHONY: swagger
swagger: ## Generate Swagger documentation
	swag init -g cmd/kubetracer/main.go -o docs/api/

# Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html

.PHONY: clean-docker
clean-docker: ## Clean Docker images
	docker rmi $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest || true

# Release targets
.PHONY: release
release: clean test build-all-platforms ## Prepare release
	@echo "Release $(VERSION) ready"

.PHONY: tag
tag: ## Create git tag
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

# Development environment
.PHONY: dev-setup
dev-setup: ## Setup development environment
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: dev-env
dev-env: ## Start development environment
	docker-compose -f deployments/docker/docker-compose.dev.yml up -d

.PHONY: dev-env-down
dev-env-down: ## Stop development environment
	docker-compose -f deployments/docker/docker-compose.dev.yml down

# Utility targets
.PHONY: version
version: ## Show version
	@echo $(VERSION)

.PHONY: info
info: ## Show build info
	@echo "App Name: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(GO_VERSION)"
