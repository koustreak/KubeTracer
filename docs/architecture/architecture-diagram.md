# KubeTracer - Simplified Architecture

```mermaid
graph TB
    subgraph "User Access"
        Browser[Web Browser]
        PortForward[kubectl port-forward]
        Browser --> PortForward
    end

    subgraph "KubeTracer Pod (Helm Deployed)"
        subgraph "Web Server :8080"
            WebUI[Web Interface]
            API[REST API]
            Static[Static Files]
        end
        
        subgraph "Lineage Engine"
            Scanner[Resource Scanner]
            Builder[Relationship Builder]
            GraphGen[Graph Generator]
        end
        
        subgraph "Storage"
            Memory[In-Memory Store]
        end
    end

    subgraph "Kubernetes Cluster"
        Pods[Pods]
        Services[Services]
        ConfigMaps[ConfigMaps]
        Secrets[Secrets]
        Deployments[Deployments]
        StatefulSets[StatefulSets]
        DaemonSets[DaemonSets]
        CRDs[Custom Resources]
    end

    %% User Flow
    PortForward --> WebUI
    WebUI --> API
    API --> Memory
    
    %% Data Collection
    Scanner --> Pods
    Scanner --> Services
    Scanner --> ConfigMaps
    Scanner --> Secrets
    Scanner --> Deployments
    Scanner --> StatefulSets
    Scanner --> DaemonSets
    Scanner --> CRDs
    
    %% Processing
    Scanner --> Builder
    Builder --> GraphGen
    GraphGen --> Memory
    
    %% Styling
    classDef userAccess fill:#e1f5fe,stroke:#01579b,color:#000
    classDef webLayer fill:#f3e5f5,stroke:#4a148c,color:#000
    classDef engine fill:#fff3e0,stroke:#e65100,color:#000
    classDef storage fill:#e8f5e8,stroke:#1b5e20,color:#000
    classDef k8s fill:#e3f2fd,stroke:#0d47a1,color:#000
    
    class Browser,PortForward userAccess
    class WebUI,API,Static webLayer
    class Scanner,Builder,GraphGen engine
    class Memory storage
    class Pods,Services,ConfigMaps,Secrets,Deployments,StatefulSets,DaemonSets,CRDs k8s
```

## Deployment Architecture

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant Helm as Helm
    participant K8s as Kubernetes
    participant Pod as KubeTracer Pod
    participant User as End User

    Dev->>Helm: helm install kubetracer
    Helm->>K8s: Deploy KubeTracer Pod
    K8s->>Pod: Start Container
    Pod->>K8s: Scan Resources
    K8s-->>Pod: Resource Data
    Pod->>Pod: Build Lineage Graph
    Pod->>Pod: Start Web Server :8080
    User->>Pod: kubectl port-forward 8080:8080
    User->>Pod: Access Web UI
```

## Key Benefits

- **🚀 Single Command Deployment**: `helm install kubetracer ./deployments/helm/kubetracer`
- **🔍 Auto Discovery**: Automatically scans all Kubernetes resources
- **📊 Simple Web UI**: Clean interface showing resource lineage
- **⚡ Fast Access**: In-memory storage for quick response times
- **🎯 Focused Scope**: Core lineage tracking without complexity
