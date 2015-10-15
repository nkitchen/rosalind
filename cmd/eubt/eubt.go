package main

import "bufio"
import "fmt"
import "log"
import "os"
import "rosalind/tree"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	taxa := []string{}
	for scanner.Scan() {
		taxa = append(taxa, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tree.DoUnrootedBinary(taxa, func(t *tree.Node) {
		fmt.Println(t)
		})
}
