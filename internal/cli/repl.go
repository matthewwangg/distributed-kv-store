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

	parts := strings.Fields(line)
	if len(parts) < 1 {
		return nil
	}

	command := parts[0]
	args := parts[1:]

	var err error

	switch command {

	case "join":
		err = HandleJoin(args, node)
	case "leave":
		err = HandleLeave(node)
	case "query":
		err = HandleQuery(args, node)
	case "help":
		HandleHelp()
	case "exit":
		HandleExit()
	default:
		fmt.Println("Unsupported command, please try again.")

	}

	return err
}
