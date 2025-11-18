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

	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Now listening on port 8080")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf("Error accepting incoming connection: %s\n", err)
			continue
		}

		go handleConnection(conn, handler)
	}
}

func handleConnection(conn net.Conn, handler *handler.Handler) {
	fmt.Println("Incoming connection:", conn.RemoteAddr().String())
	defer conn.Close()

	// Make new reader
	reader := bufio.NewReader(conn)
	for {
		// Pass reader to parser
		parsed, err := protocol.Parse(reader)
		if err != nil {
			res := fmt.Sprintf("ERR: %s\n", err)
			conn.Write([]byte(res))
			continue
		}

		// If parse succesful execute by handler
		res := handler.Exec(parsed)
		fmt.Printf("Exec returned: %#v\n", res)

		// Write response
		res = append(res, '\n')
		n, err := conn.Write(res)
		if err != nil {
			fmt.Printf("Error writing bytes: %s", err)
		}

		fmt.Println("Bytes written:", n)
	}
}
