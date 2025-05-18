package dht

import (
	"context"

	pb "github.com/matthewwangg/distributed-kv-store/proto/node"
)

func (n *Node) Join(ctx context.Context, req *pb.MembershipChangeRequest) (*pb.MembershipChangeResponse, error) {
	n.NodeState = StateRebuilding
	n.Peers[req.Id] = req.Addr

	var peerList []*pb.Peer
	for id, addr := range n.Peers {
		peerList = append(peerList, &pb.Peer{
			Id:   id,
			Addr: addr,
		})
	}

	err := n.ClientNotifyRebuild(peerList, req.Id, req.Addr, pb.Reason_JOIN)
	if err != nil {
		return nil, err
	}

	return &pb.MembershipChangeResponse{
		Peers:   peerList,
		Success: true,
	}, nil
}

func (n *Node) Leave(ctx context.Context, req *pb.MembershipChangeRequest) (*pb.MembershipChangeResponse, error) {
	n.NodeState = StateRebuilding
	delete(n.Peers, req.Id)

	var peerList []*pb.Peer
	for id, addr := range n.Peers {
		peerList = append(peerList, &pb.Peer{
			Id:   id,
			Addr: addr,
		})
	}

	err := n.ClientNotifyRebuild(peerList, req.Id, req.Addr, pb.Reason_LEAVE)
	if err != nil {
		return nil, err
	}

	return &pb.MembershipChangeResponse{
		Peers:   peerList,
		Success: true,
	}, nil
}

func (n *Node) NotifyRebuild(ctx context.Context, req *pb.RebuildRequest) (*pb.RebuildResponse, error) {
	n.NodeState = StateRebuilding
	if req.Reason == pb.Reason_JOIN {
		n.Peers[req.Id] = req.Addr
	} else if req.Reason == pb.Reason_LEAVE {
		delete(n.Peers, req.Id)
	}

	return &pb.RebuildResponse{
		Success: true,
	}, nil
}

func (n *Node) NotifyRebuildComplete(ctx context.Context, req *pb.RebuildRequest) (*pb.RebuildResponse, error) {
	n.NodeState = StateInDHT

	return &pb.RebuildResponse{
		Success: true,
	}, nil
}

func (n *Node) Store(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	n.MemoryStore[req.Key] = req.Value

	return &pb.StoreResponse{
		Success: true,
	}, nil
}
