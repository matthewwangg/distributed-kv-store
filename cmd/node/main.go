package main

import (
	"flag"
	"github.com/matthewwangg/distributed-kv-store/internal/cli"
	"log"

	"github.com/matthewwangg/distributed-kv-store/internal/dht"
)

var (
	nodeIdentifier = flag.String("id", "", "Node Identifier")
	peerAddress    = flag.String("peer-addr", "", "Peer-to-peer IP Address")
	dataDirectory  = flag.String("data-dir", "", "Data Directory")
)

func main() {
	flag.Parse()

	node := &dht.Node{
		ID:       *nodeIdentifier,
		PeerAddr: *peerAddress,
		DataDir:  *dataDirectory,
	}

	err := node.Start()
	if err != nil {
		log.Fatalf("Startup failed: %v", err)
	}

	cli.RunREPL(node)
}
