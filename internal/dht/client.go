package dht

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	utils "github.com/matthewwangg/distributed-kv-store/internal/utils"
	pb "github.com/matthewwangg/distributed-kv-store/proto/node"
)

func (n *Node) ClientJoin(joinAddr string) (map[string]string, error) {
	log.Printf("[TRACE] Sending join request to %s", joinAddr)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(joinAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("[ClientJoin] failed to create gRPC client: %w", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)

	res, err := client.Join(ctx, &pb.MembershipChangeRequest{
		Id:   n.ID,
		Addr: n.PeerAddr,
	})
	if err != nil || res.Success == false {
		return nil, fmt.Errorf("[ClientJoin] join RPC failed: %w", err)
	}

	peers := make(map[string]string)
	for _, peer := range res.GetPeers() {
		peers[peer.Id] = peer.Addr
	}

	err = n.ClientNotifyRebuildComplete(res.GetPeers(), pb.Reason_JOIN)
	if err != nil {
		return nil, fmt.Errorf("[ClientJoin] failed to notify rebuild complete: %w", err)
	}
	n.NodeState = StateInDHT
	log.Println("[TRACE] Join process completed")

	return peers, nil
}

func (n *Node) ClientLeave(neighborAddr string) error {
	log.Printf("[TRACE] Sending leave request via neighbor %s", neighborAddr)
	n.NodeState = StateFree

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(neighborAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("[ClientLeave] failed to create gRPC client: %w", err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)

	res, err := client.Leave(ctx, &pb.MembershipChangeRequest{
		Id:   n.ID,
		Addr: n.PeerAddr,
	})
	if err != nil || res.Success == false {
		return fmt.Errorf("[ClientLeave] leave RPC failed: %w", err)
	}

	delete(n.Peers, n.ID)

	if err := n.ClientStore(); err != nil {
		return fmt.Errorf("[ClientLeave] failed to redistribute keys: %w", err)
	}

	err = n.ClientNotifyRebuildComplete(res.GetPeers(), pb.Reason_LEAVE)
	if err != nil {
		return fmt.Errorf("[ClientLeave] failed to notify rebuild complete: %w", err)
	}

	log.Println("[TRACE] Leave process completed")
	return nil
}

func (n *Node) ClientNotifyRebuild(peerList []*pb.Peer, newPeerId string, newPeerAddr string, reason pb.Reason) error {
	log.Printf("[TRACE] Notifying peers of %s: %s", reason, newPeerId)
	err := n.ClientStore()
	if err != nil {
		return fmt.Errorf("[ClientNotifyRebuild] store error: %w", err)
	}

	for _, peer := range peerList {
		if peer.Id == n.ID || peer.Id == newPeerId {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := grpc.NewClient(peer.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			cancel()
			log.Printf("[WARN] Could not contact peer %s at %s: %v", peer.Id, peer.Addr, err)
			continue
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
			log.Printf("[WARN] NotifyRebuild failed for %s: %v", peer.Id, err)
			continue
		}
	}

	return nil
}

func (n *Node) ClientNotifyRebuildComplete(peers []*pb.Peer, reason pb.Reason) error {
	log.Printf("[TRACE] Sending NotifyRebuildComplete to peers for %s", reason)
	for _, peer := range peers {
		if peer.Id == n.ID {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, err := grpc.NewClient(peer.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			cancel()
			log.Printf("[WARN] Could not contact peer %s: %v", peer.Id, err)
			continue
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
			log.Printf("[WARN] NotifyRebuildComplete failed for %s: %v", peer.Id, err)
			continue
		}
	}
	return nil
}

func (n *Node) ClientStore() error {
	log.Println("[TRACE] Starting ClientStore for key redistribution")
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
			log.Printf("[WARN] Failed to connect to peer %s: %v", target, err)
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
			log.Printf("[WARN] Failed to store key %s to %s: %v", key, target, err)
			continue
		}

		toDelete = append(toDelete, key)
	}

	for _, key := range toDelete {
		delete(n.MemoryStore, key)
	}

	log.Println("[TRACE] ClientStore completed")
	return nil
}

func (n *Node) ClientGet(addr string, key string) (string, error) {
	log.Printf("[TRACE] Sending ClientGet for key %s to %s", key, addr)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", fmt.Errorf("[ClientGet] failed to connect to %s: %w", addr, err)
	}
	defer conn.Close()

	client := pb.NewNodeClient(conn)
	res, err := client.Get(ctx, &pb.GetRequest{
		Key: key,
	})
	if err != nil {
		return "", fmt.Errorf("[ClientGet] RPC failed for %s: %w", key, err)
	}
	if res.Success == false {
		return "", fmt.Errorf("[ClientGet] key %s was not found", key)
	}

	return res.Value, nil
}
