# KubeTracer - Enterprise Project Structure

## Overview
KubeTracer is a Kubernetes cluster scanning tool that extracts K8s components, establishes lineage relationships, and provides a web-based graph visualization.

## Project Structure Explanation

### Root Level Files
- `go.mod` & `go.sum` - Go module dependencies
- `Makefile` - Build automation and common tasks
- `Dockerfile` - Container image build instructions
- `LICENSE` - Project license
- `README.md` - Project documentation

### `/cmd` - Application Entry Point
- `cmd/kubetracer/` - Main web server application (combines scanner + web UI)

### `/internal` - Private application code (cannot be imported by other projects)
- `internal/api/` - HTTP REST API handlers, routes, middleware
- `internal/config/` - Configuration loading, validation, environment management
- `internal/graph/` - Graph data structures, algorithms for traversal/visualization
- `internal/lineage/` - Core business logic for establishing resource relationships
- `internal/models/` - Data structures representing K8s resources and relationships
- `internal/scanner/` - K8s cluster scanning logic, resource discovery, in-cluster client
- `internal/storage/` - Database/cache abstractions, data persistence
- `internal/utils/` - Utility functions and helpers for internal use

### `/internal/utils` - Internal Utility Functions
- `internal/utils/` - Utility functions and helpers for internal use

### `/web` - Web Application Assets
- `web/static/` - CSS, JS, images
- `web/templates/` - HTML templates
- `web/ui/` - Frontend components (React/Vue if needed)

### `/deployments` - Deployment Configurations
- `deployments/helm/kubetracer/` - Helm chart for K8s deployment
- `deployments/kubernetes/` - Raw Kubernetes manifests

### `/configs` - Configuration Files
- Application configuration templates
- Environment-specific configs

### `/test` - Testing
- `test/unit/` - Unit tests
- `test/integration/` - Integration tests

### `/docs` - Documentation
- API documentation
- User guides
- Architecture diagrams

### `/scripts` - Build and Utility Scripts
- Build scripts
- Database migration scripts
- Development helpers

### `/build` - Build Artifacts
- Docker files
- CI/CD configurations
- Binary outputs

## Key Components Architecture

### Scanner Component
- Runs as background service within main application
- Connects to Kubernetes API continuously
- Discovers all resources in cluster at regular intervals
- Extracts metadata and relationships
- Handles RBAC permissions

### Lineage Engine
- Analyzes resource dependencies
- Builds relationship graph
- Handles complex multi-hop relationships
- Supports custom resource definitions (CRDs)

### Graph Storage
- Stores resource nodes and edges
- Supports graph queries
- Handles large cluster data
- Provides caching layer

### Web UI
- Interactive graph visualization
- Resource detail views
- Search and filtering
- Real-time updates
- Integrated with scanner service

### API Layer
- RESTful API endpoints
- GraphQL support (optional)
- Authentication/authorization
- Rate limiting

## Helm Deployment Features

The Helm chart will include:
- Single deployment for combined scanner + web UI service
- ConfigMaps for configuration
- Secrets for sensitive data
- Service for web UI access
- Ingress for external web UI exposure
- RBAC permissions for cluster scanning
- PersistentVolumes for data storage
- ServiceMonitor for monitoring integration

## Enterprise Features

### Security
- RBAC integration
- TLS/SSL support
- Authentication providers
- Audit logging

### Scalability
- Horizontal pod scaling
- Resource limits and requests
- Health checks and probes
- Graceful shutdown handling

### Observability
- Prometheus metrics
- Structured logging
- Distributed tracing
- Health endpoints

### Configuration Management
- Environment-based configs
- Secret management
- Feature flags
- Hot reloading

This structure follows Go project layout standards while maintaining enterprise-level organization and separation of concerns. The single service architecture simplifies deployment while providing both background scanning and web-based visualization capabilities.
