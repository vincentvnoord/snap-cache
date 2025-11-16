package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/vincentvnoord/internal/protocol"
)

func main() {
	// cache := store.NewCache()

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Now listening on port 8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("Incoming connection:")
	defer conn.Close()
	// Read incoming data
	// Parse
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading line from connection: %s\n", err)
			return
		}

		protocol.Parse(line)
	}
}
