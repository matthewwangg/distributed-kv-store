package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	cli "github.com/matthewwangg/distributed-kv-store/internal/cli"
	dht "github.com/matthewwangg/distributed-kv-store/internal/dht"
)

var (
	nodeIdentifier = flag.String("id", "", "Node identifier (required)")
	peerAddress    = flag.String("peer-addr", "", "Peer-to-peer IP address (required)")
	joinAddress    = flag.String("join-addr", "", "Join address (optional)")
	dataDirectory  = flag.String("data-dir", "", "Data directory (optional: default is ./data/<id>)")
)

func main() {
	flag.Parse()

	if *nodeIdentifier == "" || *peerAddress == "" {
		log.Fatal("[Startup] Both --id and --peer-addr are required.")
	}

	dataDir := *dataDirectory
	if dataDir == "" {
		dataDir = filepath.Join("data", *nodeIdentifier)
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("[Startup] Failed to create data directory: %v", err)
	}

	node := &dht.Node{
		ID:       *nodeIdentifier,
		PeerAddr: *peerAddress,
		JoinAddr: *joinAddress,
		DataDir:  dataDir,
	}

	log.Printf("[TRACE] Starting node %s at %s (data: %s)", node.ID, node.PeerAddr, node.DataDir)

	if err := node.Start(); err != nil {
		log.Fatalf("[Startup] Failed to start node: %v", err)
	}

	if err := node.BootstrapJoin(); err != nil {
		log.Fatalf("[Startup] Failed to join node: %v", err)
	}

	cli.RunREPL(node)
}
