package main

import (
	"flag"
)

var (
	nodeIdentifier = flag.String("id", "", "Node Identifier")
	peerAddress    = flag.String("peer-addr", "", "Peer-to-peer IP Address")
	joinAddress    = flag.String("join-addr", "", "Join IP Address")
	dataDirectory  = flag.String("data-dir", "", "Data Directory")
)

func main() {
	flag.Parse()
}
