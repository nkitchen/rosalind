package main

import "bufio"
import "fmt"
import "rosalind/gene"
import "log"
import "os"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	line = strings.TrimSpace(line)

	n := 1
	for i := 0; i < len(line); i++ {
		c := line[i]
		a := gene.CodonsOf(c)
		if len(a) == 0 {
			log.Fatalf("Bad code: %c", c)
		}
		n = n * len(gene.CodonsOf(c)) % 1000000
	}
	n = n * len(gene.CodonsOf(gene.StopCode)) % 1000000
	fmt.Println(n)
} 
