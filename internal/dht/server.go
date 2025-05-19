package dht

import (
	"context"
	"errors"
	"fmt"
	"log"

	utils "github.com/matthewwangg/distributed-kv-store/internal/utils"
	pb "github.com/matthewwangg/distributed-kv-store/proto/node"
)

func (n *Node) Join(ctx context.Context, req *pb.MembershipChangeRequest) (*pb.MembershipChangeResponse, error) {
	log.Printf("[TRACE] Join request received from %s at %s", req.Id, req.Addr)
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
		return nil, fmt.Errorf("[Join] failed to notify rebuild: %w", err)
	}

	log.Println("[TRACE] Join rebuild notification successful")
	return &pb.MembershipChangeResponse{
		Peers:   peerList,
		Success: true,
	}, nil
}

func (n *Node) Leave(ctx context.Context, req *pb.MembershipChangeRequest) (*pb.MembershipChangeResponse, error) {
	log.Printf("[TRACE] Leave request received for %s", req.Id)
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
		return nil, fmt.Errorf("[Leave] failed to notify rebuild: %w", err)
	}

	log.Println("[TRACE] Leave rebuild notification successful")
	return &pb.MembershipChangeResponse{
		Peers:   peerList,
		Success: true,
	}, nil
}

func (n *Node) NotifyRebuild(ctx context.Context, req *pb.RebuildRequest) (*pb.RebuildResponse, error) {
	log.Printf("[TRACE] NotifyRebuild received: %s (%s)", req.Id, req.Reason)
	n.NodeState = StateRebuilding
	if req.Reason == pb.Reason_JOIN {
		n.Peers[req.Id] = req.Addr
	} else if req.Reason == pb.Reason_LEAVE {
		delete(n.Peers, req.Id)
	}

	err := n.ClientStore()
	if err != nil {
		return nil, fmt.Errorf("[NotifyRebuild] store operation failed: %w", err)
	}

	return &pb.RebuildResponse{
		Success: true,
	}, nil
}

func (n *Node) NotifyRebuildComplete(ctx context.Context, req *pb.RebuildRequest) (*pb.RebuildResponse, error) {
	log.Printf("[TRACE] Rebuild complete from %s", req.Id)
	n.NodeState = StateInDHT

	return &pb.RebuildResponse{
		Success: true,
	}, nil
}

func (n *Node) Store(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	log.Printf("[TRACE] Store key: %s", req.Key)
	n.MemoryStore[req.Key] = req.Value

	return &pb.StoreResponse{
		Success: true,
	}, nil
}

func (n *Node) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	if n.NodeState != StateInDHT {
		return nil, errors.New("[Get] node is not available to search")
	}

	log.Printf("[TRACE] Get key: %s", req.Key)
	target := utils.GetResponsiblePeer(req.Key, n.Peers)

	if target != n.PeerAddr {
		value, err := n.ClientGet(target, req.Key)
		if err != nil {
			return &pb.GetResponse{
				Success: false,
				Value:   "",
			}, nil
		}

		return &pb.GetResponse{
			Success: true,
			Value:   value,
		}, nil
	}

	value, ok := n.MemoryStore[req.Key]
	if !ok {
		return &pb.GetResponse{
			Success: false,
			Value:   "",
		}, nil
	}

	return &pb.GetResponse{
		Success: true,
		Value:   value,
	}, nil
}
