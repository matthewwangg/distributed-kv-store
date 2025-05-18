package dht

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	utils "github.com/matthewwangg/distributed-kv-store/internal/utils"
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

	err = n.ClientNotifyRebuildComplete(res.GetPeers(), pb.Reason_JOIN)
	if err != nil {
		return nil, fmt.Errorf("failed to notify rebuild complete: %w", err)
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

	err = n.ClientNotifyRebuildComplete(res.GetPeers(), pb.Reason_LEAVE)
	if err != nil {
		return fmt.Errorf("failed to notify rebuild complete: %w", err)
	}

	return nil
}

func (n *Node) ClientNotifyRebuild(peerList []*pb.Peer, newPeerId string, newPeerAddr string, reason pb.Reason) error {
	err := n.ClientStore()
	if err != nil {
		return err
	}

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

func (n *Node) ClientNotifyRebuildComplete(peers []*pb.Peer, reason pb.Reason) error {
	for _, peer := range peers {
		if peer.Id == n.ID {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, err := grpc.NewClient(peer.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			cancel()
			return fmt.Errorf("failed to create gRPC client to %s: %w", peer.Id, err)
		}

		client := pb.NewNodeClient(conn)

		res, err := client.NotifyRebuildComplete(ctx, &pb.RebuildRequest{
			Id:     n.ID,
			Addr:   n.PeerAddr,
			Reason: reason,
		})

		cancel()
		conn.Close()

		if err != nil || res.Success == false {
			return fmt.Errorf("failed to notify %s of rebuild complete: %v", peer.Id, err)
		}
	}
	return nil
}

func (n *Node) ClientStore() error {
	var toDelete []string

	for key, value := range n.MemoryStore {
		target := utils.GetResponsiblePeer(key, n.Peers)

		if target == n.PeerAddr {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			cancel()
			fmt.Printf("Failed to connect to peer %s: %v\n", target, err)
			continue
		}

		client := pb.NewNodeClient(conn)

		_, err = client.Store(ctx, &pb.StoreRequest{
			Key:   key,
			Value: value,
		})

		cancel()
		conn.Close()

		if err != nil {
			fmt.Printf("Failed to store %s to %s: %v\n", key, target, err)
			continue
		}

		toDelete = append(toDelete, key)
	}

	for _, key := range toDelete {
		delete(n.MemoryStore, key)
	}

	return nil
}
