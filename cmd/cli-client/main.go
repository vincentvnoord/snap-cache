package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vincentvnoord/snap-cache/internal/client"
)

func main() {
	// Client connects to server
	client, err := client.NewClient("localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer client.Close()

	fmt.Println("Connected to server at localhost:8080")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.ToLower(line) == "quit" {
			fmt.Println("Exiting client")
			return
		}

		// Split into command, key, value
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 2 {
			fmt.Println("Usage: SET key value or GET key")
			continue
		}
		cmd := strings.ToUpper(parts[0])
		key := parts[1]
		value := ""
		if len(parts) == 3 {
			value = parts[2]
		}

		res, err := client.SendCommand(cmd, key, []byte(value))

		if err != nil {
			fmt.Println(strings.TrimSpace(err.Error()))
			continue
		}

		fmt.Println(strings.TrimSpace(string(res)))
	}
}
