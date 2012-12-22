package main

import "bufio"
import "fmt"
import "log"
import "os"
import "rosalind/gene"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	rna := line[:len(line) - 1]

	prot, err := gene.Translate(rna)
	fmt.Println(prot)
	if err != nil {
		log.Fatal(err)
	}
}
