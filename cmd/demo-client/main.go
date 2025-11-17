package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// encodeLengthPrefixed encodes a string into the protocol:
// <length>\r\n<data>
func encodeLengthPrefixed(s string) []byte {
	data := []byte(s)
	lenStr := strconv.Itoa(len(data))
	return []byte(lenStr + "\r\n" + s)
}

func main() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server at localhost:8080")
	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

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

		// Encode command, key, value
		toSend := encodeLengthPrefixed(cmd)
		toSend = append(toSend, encodeLengthPrefixed(key)...)
		if cmd == "SET" {
			toSend = append(toSend, encodeLengthPrefixed(value)...)
		}

		// Send to server
		_, err = conn.Write(toSend)
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		// Read server response line
		resp, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading server response:", err)
			return
		}

		fmt.Println(strings.TrimSpace(resp))
	}
}
