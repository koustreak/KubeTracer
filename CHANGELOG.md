# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure setup
- Enterprise-level Go project layout
- Helm chart for Kubernetes deployment
- Web UI framework for graph visualization
- Core components architecture defined

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A

## [0.1.0] - 2025-08-16

### Added
- Initial project setup
- Basic folder structure following Go project layout standards
- Helm chart template for deployment
- Project documentation and structure guidelines
- Core architecture design for:
  - Kubernetes cluster scanning
  - Resource lineage tracking
  - Graph-based visualization
  - Web UI interface

### Project Structure
- `/cmd/kubetracer/` - Main application entry point
- `/internal/` - Private application modules
  - `scanner/` - Kubernetes resource scanning logic
  - `lineage/` - Resource relationship tracking
  - `graph/` - Graph data structures and algorithms
  - `api/` - HTTP API handlers
  - `storage/` - Data persistence layer
  - `config/` - Configuration management
  - `models/` - Data models
- `/pkg/` - Public library code
- `/web/` - Web application assets
- `/deployments/` - Helm charts and Kubernetes manifests
- `/docs/` - Documentation
- `/test/` - Unit and integration tests

### Infrastructure
- Helm chart with RBAC permissions
- ServiceAccount for cluster access
- ConfigMaps and Secrets management
- Ingress configuration for web UI
- Storage options (memory, PostgreSQL, Redis)

---

## Template for Future Releases

## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security improvements
