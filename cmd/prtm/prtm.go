package main

import "bufio"
import "fmt"
import "os"
import "rosalind/spectrum"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, _ := br.ReadString('\n')
	protein := strings.TrimSpace(line)

	w := float64(0)
	for i := 0; i < len(protein); i++ {
		w += spectrum.MonoisotopicMass[protein[i]]
	}
	fmt.Printf("%.3f\n", w)
} 
