package benchmark

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// Variables
	server := "localhost:8080"
	imagePath := "test.jpg"
	iterations := 10000

	// Disk read benchmark
	startDisk := time.Now()
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		panic(fmt.Errorf("Error reading image: %s", err))
	}

	diskReadDuration := time.Since(startDisk)

	// Connect to cache server
	conn, err := net.Dial("tcp", server)
	if err != nil {
		panic(fmt.Errorf("Error connecting to cache server: %s", err))
	}
	defer conn.Close()

	key := "benchmark-image"
	// Cache set
	cmd := encodeCommand("SET", key, imageBytes)

	startSet := time.Now()
	_, err = conn.Write(cmd)
	if err != nil {
		panic(fmt.Errorf("Error writing SET to cache server: %s", err))
	}

	setDuration := time.Since(startSet)
}

func encodeCommand(command string, key string, value []byte) []byte {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%d\r\n%s", len(command), command))
	out.WriteString(fmt.Sprintf("%d\r\n%s", len(key), key))
	out.WriteString(fmt.Sprintf("%d\r\n", len(value)))
	out.Write(value)

	return out.Bytes()
}
