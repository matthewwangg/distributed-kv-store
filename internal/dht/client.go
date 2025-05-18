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

	res, err := client.Join(ctx, &pb.MembershipChangeRequest{
		Id:   n.ID,
		Addr: n.PeerAddr,
	})
	if err != nil || res.Success == false {
		return nil, fmt.Errorf("failed to join DHT: %w", err)
	}

	peers := make(map[string]string)

	for _, peer := range res.GetPeers() {
		peers[peer.Id] = peer.Addr
	}

	for _, peerAddr := range peers {
		if peerAddr == n.PeerAddr {
			continue
		}
		err := n.ClientNotifyRebuildComplete(peerAddr, pb.Reason_JOIN)
		if err != nil {
			return nil, fmt.Errorf("failed to notify rebuild complete: %w", err)
		}
	}
	n.NodeState = StateInDHT

	return peers, nil
}

func (n *Node) ClientLeave(neighborAddr string) error {
	n.NodeState = StateFree

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(neighborAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to create gRPC client: %w", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)

	res, err := client.Leave(ctx, &pb.MembershipChangeRequest{
		Id:   n.ID,
		Addr: n.PeerAddr,
	})
	if err != nil || res.Success == false {
		return fmt.Errorf("failed to leave DHT: %w", err)
	}

	peers := make(map[string]string)

	for _, peer := range res.GetPeers() {
		peers[peer.Id] = peer.Addr
	}

	for _, peerAddr := range peers {
		if peerAddr == n.PeerAddr {
			continue
		}
		err := n.ClientNotifyRebuildComplete(peerAddr, pb.Reason_LEAVE)
		if err != nil {
			return fmt.Errorf("failed to notify rebuild complete: %w", err)
		}
	}

	return nil
}

func (n *Node) ClientNotifyRebuild(peerList []*pb.Peer, newPeerId string, newPeerAddr string, reason pb.Reason) error {
	for _, peer := range peerList {
		if peer.Id == n.ID || peer.Id == newPeerId {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := grpc.NewClient(peer.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			cancel()
			return fmt.Errorf("Failed to connect to peer %s at %s: %v\n", peer.Id, peer.Addr, err)
		}

		client := pb.NewNodeClient(conn)

		res, err := client.NotifyRebuild(ctx, &pb.RebuildRequest{
			Id:     newPeerId,
			Addr:   newPeerAddr,
			Reason: reason,
		})

		cancel()
		conn.Close()

		if err != nil || res.Success == false {
			return fmt.Errorf("Failed to notify rebuild to peer %s: %v\n", peer.Id, err)
		}
	}

	return nil
}

func (n *Node) ClientNotifyRebuildComplete(peerAddr string, reason pb.Reason) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(peerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to create gRPC client: %w", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)

	res, err := client.NotifyRebuildComplete(ctx, &pb.RebuildRequest{
		Id:     n.ID,
		Addr:   n.PeerAddr,
		Reason: reason,
	})

	if err != nil || res.Success == false {
		return fmt.Errorf("failed to join DHT: %w", err)
	}

	return nil
}
