package cli

import (
	"errors"
	"fmt"
	"log"
	"os"

	dht "github.com/matthewwangg/distributed-kv-store/internal/dht"
)

func HandleHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  join <addr>         Join the DHT at <addr>")
	fmt.Println("  leave               Leave the current DHT")
	fmt.Println("  query <addr> <key>  Query the DHT at <addr> for the given <key>")
	fmt.Println("  exit                Exit the CLI")
}

func HandleJoin(args []string, node *dht.Node) error {
	if len(args) < 1 {
		return errors.New("[HandleJoin] usage: join <addr>")
	}

	if node.NodeState != dht.StateFree {
		return errors.New("[HandleJoin] this node is already in DHT")
	}

	log.Printf("[TRACE] Joining the DHT at %s", args[0])
	peers, err := node.ClientJoin(args[0])
	if err != nil {
		return fmt.Errorf("[HandleJoin] join failed: %w", err)
	}

	node.NodeState = dht.StateInDHT

	for id, peerAddr := range peers {
		node.Peers[id] = peerAddr
	}

	log.Println("[TRACE] Join command completed")
	return nil
}

func HandleLeave(node *dht.Node) error {
	if node.NodeState == dht.StateFree {
		return errors.New("[HandleLeave] this node is already free")
	}

	log.Println("[TRACE] Leaving the DHT")
	if len(node.Peers) > 1 {
		neighbor := ""
		for _, peerAddr := range node.Peers {
			if peerAddr != node.PeerAddr {
				neighbor = peerAddr
				break
			}
		}

		err := node.ClientLeave(neighbor)
		if err != nil {
			return fmt.Errorf("[HandleLeave] leave failed: %w", err)
		}
	}

	node.NodeState = dht.StateFree
	node.Peers = map[string]string{node.ID: node.PeerAddr}

	log.Println("[TRACE] Leave command completed")
	return nil
}

func HandleQuery(args []string, node *dht.Node) error {
	if len(args) < 2 {
		return errors.New("[HandleQuery] usage: query <addr> <key>")
	}

	if node.NodeState != dht.StateFree {
		return errors.New("[HandleQuery] this node is already part of a DHT")
	}

	addr := args[0]
	key := args[1]

	log.Printf("[TRACE] Querying key %s from %s", key, addr)
	value, err := node.ClientGet(addr, key)
	if err != nil {
		return fmt.Errorf("[HandleQuery] query failed: %w", err)
	}

	fmt.Println(value)
	log.Printf("[TRACE] Query command completed with key %s matching value %s\n", key, value)
	return nil
}

func HandleExit() {
	log.Println("[TRACE] Exiting...")
	os.Exit(0)
}
