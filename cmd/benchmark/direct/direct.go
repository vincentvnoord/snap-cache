package direct

import (
	"fmt"
	"os"
	"time"

	"github.com/vincentvnoord/snap-cache/internal/cache"
)

func Run() {
	// Disk read benchmark
	startDisk := time.Now()
	imageBytes, err := os.ReadFile("test-1mb.jpg")
	if err != nil {
		panic(fmt.Errorf("Error reading image: %s", err))
	}
	diskReadDuration := time.Since(startDisk)

	// Initialize cache
	cache := cache.NewCache()

	// Set benchmark
	setStart := time.Now()
	cache.Set("benchmark-image", imageBytes)
	setDuration := time.Since(setStart)

	// Get benchmark
	getStart := time.Now()
	entry := cache.Get("benchmark-image")
	if entry == nil || string(entry.Value) != string(imageBytes) {
		panic("GET result mismatch with original image")
	}

	getDuration := time.Since(getStart)

	fmt.Printf("Disk read time: %s\n", diskReadDuration)
	fmt.Printf("SET time: %s\n", setDuration)
	fmt.Printf("GET time: %s\n", getDuration)
}
