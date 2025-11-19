package main

import (
	"flag"
	"fmt"

	"github.com/vincentvnoord/snap-cache/cmd/benchmark/direct"
	"github.com/vincentvnoord/snap-cache/cmd/benchmark/endtoend"
)

func main() {
	benchmarkType := flag.String("bench", "endtoend", "Choose benchmark: endtoend or direct")
	flag.Parse()

	switch *benchmarkType {
	case "endtoend":
		fmt.Println("Running end-to-end benchmark...")
		endtoend.Run()
	case "direct":
		fmt.Println("Running in-memory direct benchmark...")
		direct.Run()
	default:
		fmt.Println("Unknown benchmark type. Use 'endtoend' or 'direct'.")
	}
}
