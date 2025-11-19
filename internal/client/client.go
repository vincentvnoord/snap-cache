package client

import (
	"bytes"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(serverAddress string) (*Client, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}

	return &Client{conn}, nil
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

	return []byte{}, nil
}

func encodeProtocol(command string, key string, value []byte) []byte {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%d\r\n%s", len(command), command))
	out.WriteString(fmt.Sprintf("%d\r\n%s", len(key), key))
	out.WriteString(fmt.Sprintf("%d\r\n", len(value)))
	out.Write(value)

	return out.Bytes()
}
