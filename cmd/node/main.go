package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	cli "github.com/matthewwangg/distributed-kv-store/internal/cli"
	dht "github.com/matthewwangg/distributed-kv-store/internal/dht"
)

var (
	nodeIdentifier = flag.String("id", "", "Node identifier (required)")
	peerAddress    = flag.String("peer-addr", "", "Peer-to-peer IP address (required)")
	dataDirectory  = flag.String("data-dir", "", "Data directory (optional: default is ./data/<id>)")
)

func main() {
	flag.Parse()

	if *nodeIdentifier == "" || *peerAddress == "" {
		log.Fatal("Both --id and --peer-addr are required.")
	}

	dataDir := *dataDirectory
	if dataDir == "" {
		dataDir = filepath.Join("data", *nodeIdentifier)
	}

	node := &dht.Node{
		ID:       *nodeIdentifier,
		PeerAddr: *peerAddress,
		DataDir:  dataDir,
	}

	fmt.Printf("Starting node %s at %s (data: %s)\n", node.ID, node.PeerAddr, node.DataDir)

	err := node.Start()
	if err != nil {
		log.Fatalf("Startup failed: %v", err)
	}

	cli.RunREPL(node)
}
