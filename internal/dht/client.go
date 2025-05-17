package dht

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/matthewwangg/distributed-kv-store/proto/node"
)

func (n *Node) ClientJoin(joinAddr string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(joinAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)

	res, err := client.Join(ctx, &pb.JoinRequest{
		Id:   n.ID,
		Addr: n.PeerAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to join DHT: %w", err)
	}

	peers := make(map[string]string)

	for _, peer := range res.GetPeers() {
		peers[peer.Id] = peer.Addr
	}

	return peers, nil
}
