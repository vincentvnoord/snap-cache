package direct

import (
	"fmt"
	"os"
	"time"

	"github.com/vincentvnoord/snap-cache/internal/cache"
)

func Run() {
	iterations := 10000

	// Disk read benchmark
	startDisk := time.Now()
	imageBytes, err := os.ReadFile("test-1mb.jpg")
	if err != nil {
		panic(fmt.Errorf("Error reading image: %s", err))
	}
	diskReadDuration := time.Since(startDisk)

	// Initialize cache
	cache := cache.NewCache()

	var totalSet time.Duration
	for range iterations {
		// Set benchmark
		setStart := time.Now()
		cache.Set("benchmark-image", imageBytes)
		totalSet += time.Since(setStart)
	}
	avgSet := totalSet / time.Duration(iterations)

	var totalGet time.Duration
	for range iterations {
		// Get benchmark
		getStart := time.Now()
		entry := cache.Get("benchmark-image")
		if entry == nil || string(entry.Value) != string(imageBytes) {
			panic("GET result mismatch with original image")
		}

		totalGet += time.Since(getStart)
	}

	avgGet := totalGet / time.Duration(iterations)

	fmt.Printf("Disk read time: %s\n", diskReadDuration)
	fmt.Printf("Average SET time: %s\n", avgSet)
	fmt.Printf("Average GET time: %s\n", avgGet)
	fmt.Printf("Total GET time: %s\n", totalGet)
}
