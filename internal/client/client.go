package client

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"github.com/vincentvnoord/snap-cache/internal/protocol"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

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
	res, err := protocol.ReadLine(c.reader)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func encodeProtocol(command string, key string, value []byte) []byte {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%d\r\n%s", len(command), command))
	out.WriteString(fmt.Sprintf("%d\r\n%s", len(key), key))
	out.WriteString(fmt.Sprintf("%d\r\n", len(value)))
	out.Write(value)

	return out.Bytes()
}
