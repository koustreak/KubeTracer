# KubeTracer - Simplified Architecture

## Overview
KubeTracer is a Kubernetes lineage tracking tool deployed via Helm with a simple web interface on localhost.

```
                    ┌─────────────────────────────────────────────┐
                    │              USER ACCESS                    │
                    │                                             │
                    │        Browser → localhost:8080             │
                    │            (Web Interface)                  │
                    └─────────────────┬───────────────────────────┘
                                      │
                    ┌─────────────────▼───────────────────────────┐
                    │            KUBETRACER POD                   │
                    │         (Deployed via Helm)                │
                    │                                             │
                    │  ┌─────────────┐    ┌─────────────────────┐ │
                    │  │  WEB SERVER │    │   LINEAGE ENGINE    │ │
                    │  │             │    │                     │ │
                    │  │ • Static    │    │ • Resource Scanner  │ │
                    │  │   Files     │◄───┤ • Relationship     │ │
                    │  │ • Templates │    │   Builder           │ │
                    │  │ • API       │    │ • Graph Generator   │ │
                    │  └─────────────┘    └─────────────────────┘ │
                    │           │                   │             │
                    │           │                   ▼             │
                    │           │         ┌─────────────────────┐ │
                    │           │         │   IN-MEMORY STORE   │ │
                    │           │         │                     │ │
                    │           │         │ • Resource Data     │ │
                    │           └────────►│ • Relationships     │ │
                    │                     │ • Lineage Graph     │ │
                    │                     └─────────────────────┘ │
                    └─────────────────┬───────────────────────────┘
                                      │
                                      ▼
                    ┌─────────────────────────────────────────────┐
                    │           KUBERNETES CLUSTER                │
                    │                                             │
                    │  Resources Monitored:                       │
                    │  • Pods                                     │
                    │  • Services                                 │
                    │  • ConfigMaps                               │
                    │  • Secrets                                  │
                    │  • Deployments                              │
                    │  • StatefulSets                             │
                    │  • DaemonSets                               │
                    │  • Custom Resources (CRDs)                  │
                    └─────────────────────────────────────────────┘
```

## Deployment Flow

```
Developer/Ops Team
       │
       ▼
helm install kubetracer ./deployments/helm/kubetracer
       │
       ▼
KubeTracer Pod Starts
       │
       ▼ 
Scans Kubernetes Resources
       │
       ▼
Builds Lineage Graph
       │
       ▼
Web Interface Ready at localhost:8080
```

## Core Components

### 1. **Web Server** (Port 8080)
- Serves static HTML/CSS/JS files
- Provides REST API for lineage data
- Renders relationship graphs
- Simple, lightweight interface

### 2. **Lineage Engine**
- Kubernetes API client
- Resource discovery and monitoring
- Relationship analysis
- Graph generation

### 3. **In-Memory Store**
- Fast access to lineage data
- Resource relationships cache
- Real-time updates

## Key Features

- **🚀 Easy Deployment**: Single Helm command
- **🔍 Auto Discovery**: Automatically finds all K8s resources
- **📊 Visual Lineage**: Web-based dependency graphs
- **⚡ Real-time**: Live updates as resources change
- **🎯 Focused**: Core lineage tracking without complexity

## Access Pattern

```
User → kubectl port-forward → localhost:8080 → KubeTracer Web UI
```

## Minimal Dependencies
- Kubernetes cluster
- Helm 3.x
- Web browser
- kubectl (for port-forwarding)
