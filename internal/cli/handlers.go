package cli

import (
	"errors"
	"fmt"
	"os"

	dht "github.com/matthewwangg/distributed-kv-store/internal/dht"
)

func HandleHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  join <addr>   Join another node in the DHT")
	fmt.Println("  exit			 Exit the CLI")
}

func HandleJoin(args []string, node *dht.Node) error {
	if len(args) < 1 {
		return errors.New("no address provided")
	}

	if node.NodeState == dht.StateInDHT {
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

func HandleExit() {
	fmt.Println("Exiting...")
	os.Exit(0)
}
