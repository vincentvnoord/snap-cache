package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/vincentvnoord/snap-cache/internal/client"
)

func main() {
	server := "localhost:8080"
	imagePath := "test-1mb.jpg"
	iterations := 10000

	fmt.Printf("Running benchmark on %s...\n", server)

	// Disk read benchmark
	startDisk := time.Now()
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		panic(fmt.Errorf("Error reading image: %s", err))
	}
	diskReadDuration := time.Since(startDisk)

	// Connect to server
	client, err := client.NewClient(server)
	if err != nil {
		panic(fmt.Errorf("Error connecting to cache server: %s", err))
	}
	defer client.Close()

	key := "benchmark-image"

	// --- 1. SET once ---
	setStart := time.Now()
	res, err := client.SendCommand("SET", key, imageBytes)
	if err != nil {
		panic(fmt.Errorf("Error writing SET: %s", err))
	}
	if !bytes.Equal(res, []byte("OK")) {
		panic(fmt.Errorf("Expected OK, got: %s", res))
	}
	setDuration := time.Since(setStart)

	// --- 2. Benchmark GET for N iterations ---
	var totalGet time.Duration

	for range iterations {
		start := time.Now()
		res, err = client.SendCommand("GET", key, nil)
		if err != nil {
			panic(fmt.Errorf("Error sending GET: %s", err))
		}

		if !bytes.Equal(res, imageBytes) {
			panic("GET result mismatch with original image")
		}

		totalGet += time.Since(start)
	}

	avgGet := totalGet / time.Duration(iterations)

	// --- Results ---
	fmt.Println("=== Benchmark Results ===")
	fmt.Printf("Disk read:          %v\n", diskReadDuration)
	fmt.Printf("SET (1 time):       %v\n", setDuration)
	fmt.Printf("GET total (%d):     %v\n", iterations, totalGet)
	fmt.Printf("GET avg:            %v\n", avgGet)
}
