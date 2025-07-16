
# Microservice Architecture Assistant

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue)](https://golang.org/)
[![Helm Version](https://img.shields.io/badge/Helm-3+-blue)](https://helm.sh/)
[![Kubernetes Compatible](https://img.shields.io/badge/Kubernetes-1.27+-326CE5)](https://kubernetes.io)

Automated analysis and optimization system for microservice-based applications with intelligent Kubernetes deployment generation.

## ✨ Core Features
- **Architectural Analysis Engine**:
  - Static parsing of Docker configurations
  - Dependency graph construction
  - Cross-service relationship detection
- **Intelligent Service Clustering**:
  - Louvain algorithm for resource-aware grouping
  - Shared dependency identification (DBs, caches)
  - Communication pattern analysis
- **Production-Ready Outputs**:
  - Optimized Helm charts with CI/CD-ready templating
  - API Gateway and cache auto-configuration
  - Namespace isolation recommendations

## 📊 Validation Example
Tested on mvp application with 11 components (6 services + 5 DBs):
> <img width="437" height="409" alt="image" src="https://github.com/user-attachments/assets/157cd06e-b6a3-466f-86e2-1b2a6198de4b" />

Giving the following structure:
> <img width="624" height="610" alt="image" src="https://github.com/user-attachments/assets/bc517e29-78b7-464c-b3e0-c1f980995e76" />

Resulting in:
```text
✅ 4 logical clusters generated
✅ faster deployment cycles
✅ reduction in manual config work
✅ Zero-config service meshing
```

## ⚙️ Technical Implementation
| Component               | Implementation Details              |
|-------------------------|-------------------------------------|
| Core Analyzer           | Go (YAML parsing, concurrent processing) |
| Cluster Orchestration   | Kubernetes with Helm v3             |
| Optimization Algorithms | Weighted graph analysis (Louvain)   |
