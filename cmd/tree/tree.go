package main

import "fmt"
import "log"
import "os"
import "rosalind/graph"

func main() {
	g, err := graph.ReadDirectedAdjacencies(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	cc := graph.ConnectedComponents(g)
	fmt.Println(len(cc) - 1)
}
