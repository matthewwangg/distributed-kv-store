# distributed-kv-store

A lightweight distributed key-value store in Go using gRPC and consistent hashing. Nodes can join or leave a peer-to-peer cluster, and keys are rebalanced accordingly.

## ✨ Features

- Dynamic join/leave with peer notification
- Consistent hashing for key ownership
- In-memory key-value storage
- REPL-based CLI with basic commands
- Per-node logging to `logs/<node-id>/`

## 📁 Project Structure

```
distributed-kv-store
├── cmd/                    # Entry point (main.go)
├── internal/
│   ├── dht/                # DHT node logic, gRPC server/client
│   ├── cli/                # REPL + command handlers
│   └── utils/              # Logging, key hashing, store I/O
├── proto/                  # node.proto definitions
├── data/                   # Per-node key-value files (optional)
└── logs/                   # Per-node logs (auto-created)
```

## 🚀 Usage

### Run with Docker (Recommended)

#### Build the Docker Image

```bash
docker build -t kv-store-node .
```

#### Run a Node

```bash
docker run -it --rm \
  --name=node1 \
  kv-store-node \
  --id=node1 --peer-addr=node1:8080
```

> ✅ Data files should be embedded in the image under `/data/node1/` via `COPY data /data` in the Dockerfile.

#### Multi-Node Example

Create a Docker network:

```bash
docker network create kv-net
```

Start node1:

```bash
docker run -it --rm --network=kv-net \
  --name=node1 -p 8080:8080 \
  kv-store-node \
  --id=node1 --peer-addr=node1:8080
```

Start node2:

```bash
docker run -it --rm --network=kv-net \
  --name=node2 -p 8081:8081 \
  kv-store-node \
  --id=node2 --peer-addr=node2:8081
```

Join from node2's CLI:

```text
join node1:8080
```

### Run Without Docker (Optional)

#### Build Locally

```bash
go build -o kv-node ./cmd
```

#### Start a Node

```bash
./kv-node --id=node1 --peer-addr=127.0.0.1:8001
```

If `--data-dir` is not set, defaults to `./data/<id>`.

### REPL Commands

```text
join <addr>         Join the DHT via a peer node
leave               Leave the DHT and rebalance keys
query <addr> <key>  Query a specific key from a peer
help                Show available commands
exit                Exit the REPL
```
