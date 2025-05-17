package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	dht "github.com/matthewwangg/distributed-kv-store/internal/dht"
)

func RunREPL(node *dht.Node) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(". Type 'help' for commands.")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if err := Dispatch(line, node); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func Dispatch(line string, node *dht.Node) error {
	line = strings.TrimSpace(line)

	switch line {

	case "help":
		fmt.Println("help command")
	case "join":
		fmt.Println("join command")
	case "exit":
		fmt.Println("exit command")
	default:
		fmt.Println("unsupported command")

	}

	return nil
}
