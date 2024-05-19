package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Terminal error. Usage: go run main.go <inputfile>")
		return
	}
	filename := os.Args[1]
	numAnts, graph, startRoom, endRoom := input(filename)
	maxFlow := edmondsKarp(graph, startRoom, endRoom)
	fmt.Printf("Max flow: %d\n", maxFlow)
	simAnts(numAnts, graph, startRoom, endRoom)
}
