package dht

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/matthewwangg/distributed-kv-store/proto/node"
)

type Node struct {
	ID       string            `json:"id"`
	PeerAddr string            `json:"peerAddr"`
	JoinAddr string            `json:"joinAddr"`
	DataDir  string            `json:"dataDir"`
	Store    map[string]string `json:"store"`
	Peers    map[string]string `json:"peers"`

	pb.UnimplementedNodeServer
}

func (n *Node) Start() error {
	if n.JoinAddr != "" {
		peers, err := n.ClientJoin()
		if err != nil {
			log.Fatalf("Error while attempting to join DHT at %s: %v", n.JoinAddr, err)
		}
		n.Peers = peers
	} else {
		n.Peers = make(map[string]string)
	}

	n.Peers[n.ID] = n.PeerAddr

	lis, err := net.Listen("tcp", n.PeerAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", n.PeerAddr, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNodeServer(grpcServer, n)

	fmt.Printf("Node listening at %s\n", n.PeerAddr)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return nil
}
