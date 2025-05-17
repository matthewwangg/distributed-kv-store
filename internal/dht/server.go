package dht

import (
	"context"

	pb "github.com/matthewwangg/distributed-kv-store/proto/node"
)

func (n *Node) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	n.NodeState = StateRebuilding
	n.Peers[req.Id] = req.Addr

	var peerList []*pb.Peer
	for id, addr := range n.Peers {
		peerList = append(peerList, &pb.Peer{
			Id:   id,
			Addr: addr,
		})
	}

	err := n.ClientNotifyRebuild(peerList)
	if err != nil {
		return nil, err
	}

	return &pb.JoinResponse{Peers: peerList}, nil
}

func (n *Node) NotifyRebuild(ctx context.Context, req *pb.RebuildRequest) (*pb.RebuildResponse, error) {
	n.NodeState = StateRebuilding

	return &pb.RebuildResponse{Success: true}, nil
}
