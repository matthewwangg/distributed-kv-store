# distributed-kv-store

A lightweight distributed key-value store in Go using gRPC and consistent hashing. Nodes can join or leave a peer-to-peer cluster, and keys are rebalanced accordingly.

---

## ✨ Features

- Dynamic join/leave with peer notification
- Consistent hashing for key ownership
- In-memory key-value storage
- REPL-based CLI with basic commands
- Kubernetes support via StatefulSet and headless Service
- Per-node logging to `logs/<node-id>/`

---

## 📁 Project Structure

```text
distributed-kv-store
├── cmd/                    # Entry point (main.go)
├── internal/
│   ├── dht/                # DHT node logic, gRPC server/client
│   ├── cli/                # REPL + command handlers
│   └── utils/              # Logging, key hashing, store I/O
├── proto/                  # node.proto definitions
├── data/                   # Optional key-value seed files
├── k8s/                    # Kubernetes configs (kind, base manifests)
└── logs/                   # Per-node logs (auto-created)
```

---

## 💻 REPL Commands

```text
join <addr>         Join the DHT via a peer node  
leave               Leave the DHT and rebalance keys  
query <addr> <key>  Query a specific key from a peer  
help                Show available commands  
exit                Exit the REPL  
```

---

## 🚀 Deployment

### 🧩 Kubernetes (Recommended)

The project supports deploying 5+ nodes using a StatefulSet and a headless Service.

#### Prerequisites

- Docker
- kind
- kubectl

#### Create a Local Cluster

```bash
kind create cluster --name dkv --config k8s/kind/kind-config.yaml
```

#### Load the Image into kind

```bash
docker build -t kv-node:latest .
kind load docker-image kv-node:latest --name dkv
```

#### Deploy to Kubernetes

```bash
kubectl apply -f k8s/base/service.yaml
kubectl apply -f k8s/base/statefulset.yaml
```

#### View Logs

```bash
kubectl logs kv-store-0
kubectl logs kv-store-1
```

Each pod automatically joins the cluster on startup and redistributes keys.

---

### 🐳 Docker (Optional)

#### Build the Docker Image

```bash
docker build -t kv-node .
```

#### Run a Single Node

```bash
docker run -it --rm kv-node \
  --id=node1 --peer-addr=node1:8080
```

#### Multi-Node Example with Docker Network

```bash
docker network create kv-net
```

Start node1:

```bash
docker run -it --rm --network=kv-net --name=node1 \
  kv-node --id=node1 --peer-addr=node1:8080
```

Start node2:

```bash
docker run -it --rm --network=kv-net --name=node2 \
  kv-node --id=node2 --peer-addr=node2:8081
```

Then from the REPL in node2:

```text
join node1:8080
```

---

### 🧪 Run Locally (Go)

```bash
go build -o kv-node ./cmd
./kv-node --id=node1 --peer-addr=127.0.0.1:8001
```

If `--data-dir` is not set, it defaults to `./data/<id>`.

