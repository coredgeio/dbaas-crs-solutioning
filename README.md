# KubeDB Multi-Database Deployment Suite 🚀

[![KubeDB](https://img.shields.io/badge/KubeDB-Production%20Ready-blue)](https://kubedb.com)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![Author](https://img.shields.io/badge/Author-Amitesh%20Singh-orange)](https://github.com/amiteshsingh)

**Production-ready KubeDB manifests** for deploying enterprise-grade databases on Kubernetes with external access via MetalLB LoadBalancer.

This repository provides battle-tested YAML configurations for six major databases, each with custom authentication, persistent storage, and external connectivity out of the box.

---

## 📦 Databases Included

| Database | Version | Mode/Topology | External Port | Auth Secret |
|----------|---------|---------------|---------------|-------------|
| **MongoDB** | 6.0.12 | ReplicaSet (rs0) | 27017 | mongo-custom-auth |
| **PostgreSQL** | 16.10 | Hot Standby | 5432 | pg-custom-auth |
| **MySQL** | 8.0.35 | Standalone | 3306 | mysql-custom-auth |
| **Redis** | 8.2.2 | Standalone | 6379 | redis-custom-auth |
| **Kafka** | 3.7.2 | KRaft (Controller+Broker) | 9092, 9093 | kafka-custom-auth |
| **ClickHouse** | 25.7.1 | Standalone | 8123 (HTTP), 9000 (Native) | ch-custom-auth |

---

## 🎯 Features

✅ **Custom Authentication** - Pre-configured users and passwords (externally managed)  
✅ **MetalLB Integration** - Automatic LoadBalancer IP assignment from dedicated pools  
✅ **Persistent Storage** - Uses `local-path` StorageClass (configurable)  
✅ **Production Topology** - ReplicaSets, Hot Standby, KRaft mode where applicable  
✅ **GitOps Ready** - Declarative YAML manifests for version control  
✅ **Namespace Isolation** - Each database in its own namespace  

---

## 🔧 Prerequisites

### 1. Kubernetes Cluster
- **Version**: v1.25+
- **CNI**: Functional network plugin (Calico, Flannel, etc.)
- **StorageClass**: `local-path` (or modify manifests for your provider)

### 2. KubeDB Operator

Install KubeDB with required database support:

```bash
helm upgrade --install kubedb oci://ghcr.io/appscode-charts/kubedb \
  --version v2025.10.17 \
  --namespace kubedb \
  --create-namespace \
  --set global.featureGates.ClickHouse=true \
  --set-file global.license=/path/to/kubedb-license.txt
```

Verify installation:

```bash
kubectl get pods -n kubedb
kubectl get crds | grep kubedb.com
```

### 3. MetalLB LoadBalancer

MetalLB must be installed and configured with IP address pools.

**Example IPAddressPool configurations** referenced in manifests:

```yaml
---
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: sandbox
  namespace: metallb-system
spec:
  addresses:
  - 172.16.109.207-172.16.109.209
```

Apply MetalLB configuration:

```bash
kubectl apply -f metallb-config.yaml
```

---

## 📁 Repository Structure

```
dbs/
├── mongoall.yaml          # MongoDB ReplicaSet with LoadBalancer
├── postgresall.yaml       # PostgreSQL Hot Standby with LoadBalancer
├── mysqlall.yaml          # MySQL Standalone with LoadBalancer
├── redisall.yaml          # Redis Standalone with LoadBalancer
├── kafkaall.yaml          # Kafka KRaft mode with separate broker/controller services
├── clickhouseall.yaml     # ClickHouse with HTTP and Native protocol access
└── README.md              # This file
```

---

## 🚀 Deployment Guide

### MongoDB

**Features:**
- ReplicaSet (rs0) with 1 replica
- Custom auth via `MongoDBOpsRequest` (RotateAuth)
- LoadBalancer for primary and standby services
- MetalLB pool: `sandbox`

**Deploy:**

```bash
kubectl create namespace mongo
kubectl apply -f mongoall.yaml
```

**Verify:**

```bash
kubectl get mongodb -n mongo
kubectl get svc -n mongo
kubectl get secret mongo-custom-auth -n mongo -o yaml
```

**Connect:**

```bash
# Get LoadBalancer IP
LB_IP=$(kubectl get svc mongo-rs -n mongo -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# Connect with mongosh
mongosh "mongodb://root:admin123@${LB_IP}:27017"

---

### PostgreSQL

**Features:**
- Hot Standby with streaming replication
- Custom auth secret (externally managed)
- Leader election configured
- MetalLB pool: `newip`

**Deploy:**

```bash
kubectl create namespace postgres
kubectl apply -f postgresall.yaml
```

**Verify:**

```bash
kubectl get postgres -n postgres
kubectl get svc -n postgres
```

**Connect:**

```bash
# Get LoadBalancer IP
LB_IP=$(kubectl get svc pg-cluster-primary -n postgres -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# Connect with psql
psql "postgresql://postgres:postgres123@${LB_IP}:5432/postgres"

---

### MySQL

**Features:**
- Standalone deployment with 1 replica
- Externally managed auth secret
- Durable storage with local-path
- MetalLB pool: `choose yours`

**Deploy:**

```bash
kubectl create namespace mysql
kubectl apply -f mysqlall.yaml
```

**Verify:**

```bash
kubectl get mysql -n mysql
kubectl get svc -n mysql
```

**Connect:**

```bash
# Get LoadBalancer IP
LB_IP=$(kubectl get svc mysql-cluster-primary -n mysql -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# Connect with mysql client
mysql -h ${LB_IP} -P 3306 -u root -pmysql123

---

### Redis

**Features:**
- Standalone mode with AUTH enabled
- Custom authentication
- Durable persistent storage
- MetalLB pool: `cinderip`

**Deploy:**

```bash
kubectl create namespace redis
kubectl apply -f redisall.yaml
```

**Verify:**

```bash
kubectl get redis -n redis
kubectl get svc -n redis
```

**Connect:**

```bash
# Get LoadBalancer IP
LB_IP=$(kubectl get svc redis-cluster-primary -n redis -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# Connect with redis-cli
redis-cli -h ${LB_IP} -p 6379 -a redis123


---

### Kafka

**Features:**
- KRaft mode (no ZooKeeper required)
- Separate controller and broker services
- SASL/PLAIN authentication
- MetalLB pool: `choose yours`

**Deploy:**

```bash
kubectl create namespace kafka
kubectl apply -f kafkaall.yaml
```

**Verify:**

```bash
kubectl get kafka -n kafka
kubectl get svc -n kafka
```

**Connect:**

```bash
# Get LoadBalancer IPs
BROKER_IP=$(kubectl get svc kafka-broker-external -n kafka -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
CONTROLLER_IP=$(kubectl get svc kafka-controller-external -n kafka -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# Create topic (using kafka-topics.sh)
kafka-topics.sh --bootstrap-server ${BROKER_IP}:9092 \
  --command-config <(echo "security.protocol=SASL_PLAINTEXT
sasl.mechanism=PLAIN
sasl.jaas.config=org.apache.kafka.common.security.plain.PlainLoginModule required username='' password='';") \
  --create --topic test-topic --partitions 3 --replication-factor 1

**Exposed Services:**
- Broker: Port 9092
- Controller: Port 9093

---

### ClickHouse

**Features:**
- Standalone deployment
- HTTP interface (8123) and Native protocol (9000)
- Custom authentication
- MetalLB pool: `sandbox`

**Deploy:**

```bash
kubectl create namespace clickhouse
kubectl apply -f clickhouseall.yaml
```

**Verify:**

```bash
kubectl get clickhouse -n clickhouse
kubectl get svc -n clickhouse
```

**Connect:**

```bash
# Get LoadBalancer IP
LB_IP=$(kubectl get svc ch-cluster-external -n clickhouse -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# HTTP Interface
curl -u admin:clickhouse123 "http://${LB_IP}:8123/?query=SELECT%20version()"

# Native Protocol (clickhouse-client)
clickhouse-client --host ${LB_IP} --port 9000 --user admin --password clickhouse123


**Exposed Ports:**
- 8123: HTTP interface
- 9000: Native protocol

---

## 🔒 Security Considerations

### Production Recommendations

1. **Change Default Passwords** - Update auth secrets before deployment:
```bash
echo -n "your-secure-password" | base64
```

2. **Enable TLS/SSL** - Add TLS configuration to KubeDB specs:
```yaml
spec:
  enableSSL: true
  tls:
    issuerRef:
      name: ca-issuer
      kind: Issuer
      apiGroup: cert-manager.io
```

3. **Network Policies** - Restrict access to databases:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: db-network-policy
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: kafkas.kubedb.com
  policyTypes:
  - Ingress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: allowed-namespace
```

4. **RBAC** - Implement least-privilege access controls

5. **Backup Strategy** - Configure automated backups using KubeDB Stash integration

---

## 🛠️ Troubleshooting

### Database Not Starting

```bash
# Check pod status
kubectl get pods -n <namespace>

# View pod logs
kubectl logs <pod-name> -n <namespace>

# Describe database resource
kubectl describe <db-type> <db-name> -n <namespace>
```

### LoadBalancer IP Pending

```bash
# Check MetalLB configuration
kubectl get ipaddresspool -n metallb-system
kubectl logs -n metallb-system -l app=metallb

# Verify service annotation
kubectl get svc <service-name> -n <namespace> -o yaml | grep metallb
```

### Storage Issues

```bash
# Check PVC status
kubectl get pvc -n <namespace>

# Describe PVC
kubectl describe pvc <pvc-name> -n <namespace>

# Verify StorageClass
kubectl get storageclass local-path -o yaml
```

### Authentication Failures

```bash
# Verify secret exists
kubectl get secret <secret-name> -n <namespace>

# Decode secret
kubectl get secret <secret-name> -n <namespace> -o jsonpath='{.data.password}' | base64 -d
```

---

## 🎓 Advanced Configuration

### Custom StorageClass

Replace `local-path` in manifests with your StorageClass:

```yaml
storage:
  storageClassName: csi-cinder-sc-retain  # OpenStack Cinder
  # storageClassName: ebs-sc               # AWS EBS
  # storageClassName: azuredisk-sc         # Azure Disk
```

### Scaling ReplicaSets

Increase replicas for high availability:

```yaml
spec:
  replicas: 3  # For MongoDB/MySQL
  
  topology:    # For Kafka
    broker:
      replicas: 3
    controller:
      replicas: 3
```

### Resource Limits

Add resource constraints:

```yaml
spec:
  podTemplate:
    spec:
      resources:
        requests:
          cpu: "1"
          memory: 2Gi
        limits:
          cpu: "2"
          memory: 4Gi
```

---

## 📊 Monitoring & Observability

### Enable Prometheus Monitoring

Add monitoring configuration to KubeDB specs:

```yaml
spec:
  monitor:
    agent: prometheus.io/operator
    prometheus:
      serviceMonitor:
        labels:
          release: prometheus
```

### Grafana Dashboards

KubeDB provides pre-built Grafana dashboards for each database type. Import them from:
- [KubeDB Grafana Dashboards](https://github.com/kubedb/installer/tree/master/charts/kubedb-metrics/dashboards)

---

## 🧪 Testing Deployments

### Automated Health Checks

```bash
#!/bin/bash
# test-deployments.sh

DATABASES=("mongo" "postgres" "mysql" "redis" "kafka" "clickhouse")

for db in "${DATABASES[@]}"; do
  echo "Testing $db..."
  kubectl wait --for=condition=Ready pods -l app.kubernetes.io/instance=$db-cluster -n $db --timeout=300s
  
  if [ $? -eq 0 ]; then
    echo "✅ $db is ready"
  else
    echo "❌ $db failed to start"
  fi
done
```

---

## 📚 Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│              Kubernetes Cluster                     │
│                                                     │
│  ┌─────────────┐  ┌─────────────┐  ┌────────────┐ │
│  │  MongoDB    │  │ PostgreSQL  │  │   MySQL    │ │
│  │  Namespace  │  │  Namespace  │  │ Namespace  │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬─────┘ │
│         │                │                 │        │
│  ┌──────┴──────┐  ┌──────┴──────┐  ┌──────┴─────┐ │
│  │  Redis      │  │  Kafka      │  │ ClickHouse │ │
│  │  Namespace  │  │  Namespace  │  │ Namespace  │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬─────┘ │
│         │                │                 │        │
│         └────────────────┴─────────────────┘        │
│                          ▼                          │
│                   ┌─────────────┐                   │
│                   │   MetalLB   │                   │
│                   │ LoadBalancer│                   │
│                   └──────┬──────┘                   │
└──────────────────────────┼──────────────────────────┘
                           ▼
                  External Clients
           (mongosh, psql, mysql, redis-cli, etc.)
```

---

## ✅ Final Verdict

### What This Repository Provides

✔ **Production-Safe Deployments** - Tested configurations for 6 major databases  
✔ **MetalLB Integration** - Automatic external IP assignment from dedicated pools  
✔ **KubeDB Native** - Fully managed by KubeDB operator (no manual intervention)  
✔ **Custom Authentication** - Pre-configured users with externally managed secrets  
✔ **Persistent Storage** - Durable data storage with configurable StorageClasses  
✔ **GitOps Ready** - Version-controlled YAML manifests  

### What's Not Included (But Easy to Add)

- ⚠️ TLS/SSL encryption (requires cert-manager)
- ⚠️ Automated backups (requires KubeDB Stash)
- ⚠️ Monitoring/Alerting (requires Prometheus/Grafana)
- ⚠️ Ingress controllers (alternative to LoadBalancer)

---

## 🎯 Next Steps

Want to enhance your deployment? Consider:

- **Add TLS Encryption** - Secure all database connections with cert-manager
- **Implement Backups** - Use Velero or KubeDB Stash for automated backups
- **Set Up Monitoring** - Deploy Prometheus + Grafana for observability
- **Configure Ingress** - Use NGINX/Traefik ingress instead of LoadBalancer
- **Enable HA** - Scale to multiple replicas with anti-affinity rules
- **Implement DR** - Set up cross-cluster replication for disaster recovery

---

## 📖 References

- [KubeDB Documentation](https://kubedb.com/docs/)
- [MetalLB Documentation](https://metallb.universe.tf/)
- [Kubernetes Storage Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/)
- [KubeDB GitHub](https://github.com/kubedb)

---

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](../LICENSE) file for details.

---

## 👤 Author

**Amitesh Singh**

- GitHub: [@amiteshsingh](https://github.com/amiteshsingh)
- Email: amiteshhsingh@gmail.com

---

## 🙏 Acknowledgments

- KubeDB team for the amazing database operator
- MetalLB community for bare-metal load balancing
- Kubernetes community for the robust platform

---

## 💬 Support

Need help? Open an issue or reach out:

- GitHub Issues: [Create an issue](https://github.com/amiteshsingh/kubedb-amitesh/issues)
- Documentation: Check the [KubeDB docs](https://kubedb.com/docs/)
- Community: Join the [KubeDB Slack](https://kubernetes.slack.com/messages/kubedb)

---

**Made with ❤️ for the Kubernetes community**

Just say the word if you need:
- Automated backup configurations 💾
- TLS/SSL setup guides 🔒
- Monitoring dashboards 📊
- HA/DR strategies 🚨

🚀 **Happy Database Deploying!**


