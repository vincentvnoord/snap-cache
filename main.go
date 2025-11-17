package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/vincentvnoord/internal/cache"
	"github.com/vincentvnoord/internal/handler"
	"github.com/vincentvnoord/internal/protocol"
)

func main() {
	cache := cache.NewCache()
	handler := &handler.Handler{Cache: cache}

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

		go handleConnection(conn, handler)
	}
}

func handleConnection(conn net.Conn, handler *handler.Handler) {
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

		parsed, err := protocol.Parse(line)
		if err != nil {
			conn.Write([]byte("ERR: Malformed input\n"))
			continue
		}

		res := handler.Exec(parsed)
		fmt.Printf("Exec returned: %#v\n", res)

		res = append(res, '\n')
		n, err := conn.Write(res)
		if err != nil {
			fmt.Printf("Error writing bytes: %s", err)
		}

		fmt.Println("Bytes written:", n)
	}
}
