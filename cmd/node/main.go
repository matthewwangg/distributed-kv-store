package main

import (
	"flag"

	"github.com/matthewwangg/distributed-kv-store/internal/dht"
)

var (
	nodeIdentifier = flag.String("id", "", "Node Identifier")
	peerAddress    = flag.String("peer-addr", "", "Peer-to-peer IP Address")
	joinAddress    = flag.String("join-addr", "", "Join IP Address")
	dataDirectory  = flag.String("data-dir", "", "Data Directory")
)

func main() {
	flag.Parse()

	node := &dht.Node{
		ID:       *nodeIdentifier,
		PeerAddr: *peerAddress,
		JoinAddr: *joinAddress,
		DataDir:  *dataDirectory,
	}

	node.Start()
}
