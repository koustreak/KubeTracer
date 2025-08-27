# KubeTracer

KubeTracer is a tool designed to establish and visualize the lineage between all Kubernetes components, such as ConfigMap, Pod, Service, Secret, CRD, DaemonSet, StatefulSet, and more. It helps users understand the relationships and dependencies within their Kubernetes clusters, making it easier to troubleshoot, audit, and optimize resource usage.

## Features
- Discover and map relationships between Kubernetes resources
- Visualize dependencies and lineage across the cluster
- Support for core and custom resources (CRDs)
- Useful for debugging, auditing, and compliance

## Getting Started

### Prerequisites
- Go 1.22 or higher
- Docker (for containerized deployment)
- Kubernetes cluster access
- kubectl configured

1. **Clone the repository:**
   ```sh
   git clone https://github.com/koustreak/kubetracer.git
   cd kubetracer
   ```
2. **Build the application:**
   ```sh
   make build
   ```
3. **Run KubeTracer:**
   ```sh
   ./bin/kubetracer
   ```

### Docker Deployment
1. **Build Docker image:**
   ```sh
   make docker-build
   ```
2. **Deploy to Kubernetes:**
   ```sh
   make deploy-kind  # For KIND cluster
   # or
   kubectl apply -f deployments/kubernetes/
   ```

## Supported Resources
- Pod
- Service
- ConfigMap
- Secret
- CRD (Custom Resource Definitions)
- DaemonSet
- StatefulSet
- (and more)

## Use Cases
- Visualizing application architecture
- Troubleshooting resource issues
- Auditing and compliance
- Understanding cluster configuration

## Contributing
Contributions are welcome! Please open issues or submit pull requests.

## License
This project is licensed under the MIT License.
