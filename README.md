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

### Build

```bash
go build -o kv-node ./cmd
```

### Start a Node

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
