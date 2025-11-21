package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/vincentvnoord/snap-cache/internal/protocol"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

// Usage: client, err := client.NewClient("localhost:8080")
// Remember to call client.Close() when done calling commands
func NewClient(serverAddress string) (*Client, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn,
		bufio.NewReader(conn),
	}, nil
}

// Creator of client owns lifetime and has to call close
func (c Client) Close() error {
	return c.conn.Close()
}

func (c Client) SendCommand(command string, key string, value []byte) ([]byte, error) {
	// Encode into protocol
	encoded := encodeProtocol(command, key, value)

	// Write to conn
	_, err := c.conn.Write(encoded)
	if err != nil {
		return nil, err
	}

	// Read response until \r\n
	buf, err := protocol.ReadLine(c.reader)
	if err != nil {
		return nil, err
	}

	// Response size
	size, err := strconv.Atoi(string(buf))
	if err != nil {
		return nil, err
	}

	// Read response in buf with the size
	res := make([]byte, size)
	_, err = io.ReadFull(c.reader, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func encodeProtocol(command string, key string, value []byte) []byte {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%d\r\n%s", len(command), command))
	out.WriteString(fmt.Sprintf("%d\r\n%s", len(key), key))

	if value != nil {
		out.WriteString(fmt.Sprintf("%d\r\n", len(value)))
		out.Write(value)
	}

	return out.Bytes()
}
