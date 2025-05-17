package cli

import (
	"fmt"
	"github.com/matthewwangg/distributed-kv-store/internal/dht"
	"os"
)

func HandleHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  join <addr>   Join another node in the DHT")
	fmt.Println("  exit			 Exit the CLI")
}

func HandleJoin(args []string, node *dht.Node) error {
	fmt.Println("Joining the DHT")

	return nil
}

func HandleExit() {
	fmt.Println("Exiting...")
	os.Exit(0)
}
