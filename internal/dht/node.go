package dht

import (
	"log"
	"net"

	"google.golang.org/grpc"

	utils "github.com/matthewwangg/distributed-kv-store/internal/utils"
	pb "github.com/matthewwangg/distributed-kv-store/proto/node"
)

type NodeState int

const (
	StateFree       NodeState = iota
	StateInDHT      NodeState = iota
	StateRebuilding NodeState = iota
)

type Node struct {
	ID          string            `json:"id"`
	PeerAddr    string            `json:"peerAddr"`
	DataDir     string            `json:"dataDir"`
	MemoryStore map[string]string `json:"memoryStore"`
	Peers       map[string]string `json:"peers"`
	NodeState   NodeState         `json:"nodeState"`

	pb.UnimplementedNodeServer
}

func (n *Node) Start() error {
	if err := utils.SetupLogger(n.ID, n.PeerAddr); err != nil {
		log.Fatalf("failed to set up logger: %v", err)
	}

	n.Peers = make(map[string]string)
	n.Peers[n.ID] = n.PeerAddr
	n.NodeState = StateFree

	kv, err := utils.LoadKeyValueDir(n.DataDir)
	if err != nil {
		return err
	}
	n.MemoryStore = kv

	lis, err := net.Listen("tcp", n.PeerAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", n.PeerAddr, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNodeServer(grpcServer, n)

	log.Printf("Node listening at %s", n.PeerAddr)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return nil
}
