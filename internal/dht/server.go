package dht

import (
	"context"

	pb "github.com/matthewwangg/distributed-kv-store/proto/node"
)

func (n *Node) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	n.Peers[req.Id] = req.Addr

	var peerList []*pb.Peer
	for id, addr := range n.Peers {
		peerList = append(peerList, &pb.Peer{
			Id:   id,
			Addr: addr,
		})
	}

	return &pb.JoinResponse{Peers: peerList}, nil
}
