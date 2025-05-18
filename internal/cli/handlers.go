package cli

import (
	"errors"
	"fmt"
	"os"

	dht "github.com/matthewwangg/distributed-kv-store/internal/dht"
)

func HandleHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  join <addr>   Join the DHT at <addr>")
	fmt.Println("  leave 		 Leave the current DHT")
	fmt.Println("  exit			 Exit the CLI")
}

func HandleJoin(args []string, node *dht.Node) error {
	if len(args) < 1 {
		return errors.New("no address provided")
	}

	if node.NodeState != dht.StateFree {
		return errors.New("this node is already in DHT")
	}

	fmt.Println("Joining the DHT")
	peers, err := node.ClientJoin(args[0])
	if err != nil {
		return err
	}

	node.NodeState = dht.StateInDHT

	for id, peerAddr := range peers {
		node.Peers[id] = peerAddr
	}

	return nil
}

func HandleLeave(node *dht.Node) error {
	if node.NodeState == dht.StateFree {
		return errors.New("this node is already free")
	}

	fmt.Println("Leaving the DHT")
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
			return err
		}
	}

	node.NodeState = dht.StateFree

	node.Peers = map[string]string{node.ID: node.PeerAddr}

	return nil
}

func HandleExit() {
	fmt.Println("Exiting...")
	os.Exit(0)
}
