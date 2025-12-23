# DBaaS Monitoring API

Go-based API to fetch metrics from Prometheus running in Kubernetes cluster and serve them to the frontend.

## Features

- 🚀 Fetch data from Prometheus (Kubernetes cluster)
- 📊 Pre-built queries for CPU, Memory, Disk, Pods, Nodes, and Container metrics
- 🔍 Custom query support
- 🌐 CORS enabled for frontend integration
- ⚡ Fast and lightweight
- 🐳 Docker support
- ☸️ Kubernetes deployment ready

## Architecture

```
├── main.go                 # Application entry point
├── config/
│   └── config.go          # Configuration management
├── models/
│   └── models.go          # Data models and structs
├── services/
│   └── prometheus.go      # Prometheus client and business logic
├── handlers/
│   └── handlers.go        # HTTP request handlers
├── router/
│   └── router.go          # Route definitions
├── k8s/
│   └── deployment.yaml    # Kubernetes deployment manifest
└── Dockerfile             # Docker image build configuration
```

## Prerequisites

- Go 1.21 or higher
- Access to Kubernetes cluster with Prometheus
- kubectl configured

## Kubernetes Prometheus Service

Your Prometheus service details:
```
Service: kube-prom-stack-kube-prome-prometheus
Namespace: monitoring
NodePort: 9090:31492
```

## Quick Start (Auto-Discovery)

The API **automatically discovers** your Prometheus endpoint! No manual configuration needed.

### Single Command Start:

```bash
cd /Users/amiteshsingh/Desktop/dbaaas-monitoring
go mod download
go run main.go
```

**That's it!** The API will:
1. ✅ Detect if running inside/outside Kubernetes cluster
2. ✅ Auto-discover Node IP using kubectl
3. ✅ Auto-discover NodePort from service
4. ✅ Build Prometheus URL: `http://<NODE_IP>:31492`
5. ✅ Validate connection before starting

### Setup Options

### Option 1: Run Locally (Auto-Discovery)

```bash
# Just run it - auto-discovery will handle everything
go mod download
go run main.go
```

Output will show:
```
🔍 Discovering Prometheus endpoint...
✓ Discovered Prometheus endpoint:
  - Node IP: 192.168.1.100
  - NodePort: 31492
  - Full URL: http://192.168.1.100:31492
🔗 Validating Prometheus connection...
✓ Prometheus health check passed
✅ Configuration loaded successfully
```

### Option 2: Manual Override (Optional)

If auto-discovery fails or you want to override:

```bash
export PROMETHEUS_URL=http://your-node-ip:31492
go run main.go
```

### Option 3: Deploy to Kubernetes Cluster

```bash
# Build Docker image
docker build -t dbaas-monitoring-api:latest .

# Deploy to Kubernetes (uses in-cluster service discovery)
kubectl apply -f k8s/deployment.yaml

# Access via NodePort 30080
# API: http://<NODE_IP>:30080/api/health
```

## API Endpoints

### Health Check
```
GET /api/health
```

### System Metrics

**CPU Metrics**
```
GET /api/metrics/cpu
```

**Memory Metrics**
```
GET /api/metrics/memory
```

**Disk Metrics**
```
GET /api/metrics/disk
```

### Kubernetes Metrics

**Pod Information**
```
GET /api/metrics/pods
```

**Node Information**
```
GET /api/metrics/nodes
```

**Container CPU Usage**
```
GET /api/metrics/container/cpu
```

**Container Memory Usage**
```
GET /api/metrics/container/memory
```

**All Targets Status**
```
GET /api/metrics/targets
```

### Custom Queries

**Single Query**
```
GET /api/metrics/query?query=up
```

Example:
```bash
curl "http://localhost:8080/api/metrics/query?query=up"
curl "http://localhost:8080/api/metrics/query?query=kube_pod_status_phase"
```

**Query Range (Time Series)**
```
GET /api/metrics/query_range?query=up&start=2024-01-01T00:00:00Z&end=2024-01-01T01:00:00Z
```

## How Auto-Discovery Works

The API uses intelligent discovery to find Prometheus:

1. **Check if inside Kubernetes cluster:**
   - Looks for service account token at `/var/run/secrets/kubernetes.io/serviceaccount/token`
   - If found, uses in-cluster DNS: `kube-prom-stack-kube-prome-prometheus.monitoring.svc.cluster.local:9090`

2. **If outside cluster:**
   - Runs: `kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}'`
   - Runs: `kubectl get svc -n monitoring kube-prom-stack-kube-prome-prometheus -o jsonpath='{.spec.ports[0].nodePort}'`
   - Builds URL: `http://<NODE_IP>:<NODEPORT>`

3. **Health Check:**
   - Tests connection to `/api/v1/query` endpoint
   - Verifies Prometheus is responsive

4. **Fallback:**
   - Uses `PROMETHEUS_URL` environment variable if set
   - Falls back to `http://localhost:9090` as last resort

### Response Format

All endpoints return JSON:

```json
{
  "status": "success",
  "data": {
    // Prometheus response data
  }
}
```

Error response:
```json
{
  "status": "error",
  "error": "error message"
}
```

## Common Prometheus Queries

You can use these with `/api/metrics/query`:

### Node Metrics
- **CPU Usage**: `100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`
- **Memory Usage**: `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`
- **Disk Usage**: `(1 - (node_filesystem_avail_bytes / node_filesystem_size_bytes)) * 100`
- **Network Traffic**: `rate(node_network_receive_bytes_total[5m])`

### Kubernetes Metrics
- **Pod Count**: `count(kube_pod_info)`
- **Pod Status**: `kube_pod_status_phase`
- **Node Status**: `kube_node_status_condition`
- **Container Restarts**: `kube_pod_container_status_restarts_total`
- **Deployments Ready**: `kube_deployment_status_replicas_ready`
- **PVC Usage**: `kubelet_volume_stats_used_bytes / kubelet_volume_stats_capacity_bytes * 100`

### Container Metrics
- **Container CPU**: `sum(rate(container_cpu_usage_seconds_total{container!=""}[5m])) by (pod, namespace)`
- **Container Memory**: `sum(container_memory_working_set_bytes{container!=""}) by (pod, namespace)`
- **Container Network**: `sum(rate(container_network_receive_bytes_total[5m])) by (pod)`

## Environment Variables

All variables are **optional** - auto-discovery handles everything!

- `PROMETHEUS_URL`: (Optional) Override auto-discovered Prometheus URL
- `PORT`: API server port (default: 8080)
- `ENVIRONMENT`: Environment name (default: development)

## Build for Production

```bash
# Binary build
go build -o dbaaas-api main.go
./dbaaas-api

# Docker build
docker build -t dbaas-monitoring-api:latest .
docker run -p 8080:8080 -e PROMETHEUS_URL=http://prometheus:9090 dbaas-monitoring-api
```

## Testing the API

```bash
# Health check
curl http://localhost:8080/api/health

# Get CPU metrics
curl http://localhost:8080/api/metrics/cpu

# Get all pods
curl http://localhost:8080/api/metrics/pods

# Custom query
curl "http://localhost:8080/api/metrics/query?query=up"

# Container CPU usage
curl http://localhost:8080/api/metrics/container/cpu
```

