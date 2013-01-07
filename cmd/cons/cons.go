package main

import "bufio"
import "fmt"
import "io"
import "log"
import "os"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	bases := ([]byte)("ACGT")
	baseIndex := map[byte]int{}
	for i, b := range bases {
		baseIndex[b] = i
	}

	var profile [][]int
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		dna := strings.TrimSpace(line)
		if len(profile) == 0 {
			k := len(bases)
			m := make([]int, len(dna) * k)
			profile = make([][]int, len(dna))
			for i := range dna {
				profile[i] = m[i*k:(i+1)*k]
			}
		}

		for i := 0; i < len(dna); i++ {
			profile[i][baseIndex[dna[i]]]++
		}
	}

	consensus := []byte{}
	for i := range profile {
		maxBase := byte('X')
		maxCount := -1
		for b, j := range baseIndex {
			n := profile[i][j]
			if n > maxCount {
				maxBase = b
				maxCount = n
			}
		}
		consensus = append(consensus, maxBase)
	}

	fmt.Println(string(consensus))
	for i, b := range bases {
		fmt.Printf("%c:", b)
		for j := 0; j < len(profile); j++ {
			n := profile[j][i]
			fmt.Print(" ", n)
		}
		fmt.Println()
	}
} 
