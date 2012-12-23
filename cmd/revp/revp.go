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
	dna := line[:len(line) - 1]

	for i := 1; i < len(dna) - 1; i++ {
		// Look for the center of a reverse palindrome first.
		var k int
		for k = 0; i + k < len(dna) - 1 && k < 12; k++ {
			j := i - k
			if j < 0 {
				break
			}
			if dna[i + 1 + k] != gene.DnaComplement[dna[j]] {
				break
			}
		}
		for m := k; m >= 2; m-- {
			if m > 6 {
				panic(fmt.Sprint("Unexpected length", k*2))
			}
			fmt.Println(i + 2 - m, 2 * m)
		}
	}
}
