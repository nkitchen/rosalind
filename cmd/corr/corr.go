package main

import "bufio"
import "fmt"
import "io"
import "log"
import "os"
import "rosalind/gene"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	reads := []string{}
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		reads = append(reads, strings.TrimSpace(line))
	}

	seen := map[string]int{}
	for _, s := range reads {
		seen[s] += 1
		rc := gene.ReverseComplement(s)
		if rc != s {
			seen[rc] += 1
		}
	}

	bases := ([]byte)("ACGT")
	for _, s := range reads {
		if seen[s] == 1 {
			corrected := ""
			for i := 0; i < len(s); i++ {
				for _, b := range bases {
					if s[i] != b {
						t := s[:i] + string(b) + s[i+1:]
						if seen[t] > 1 {
							if corrected == "" {
								corrected = t
							} else {
								log.Fatal("Found multiple corrections: " + s)
							}
						}
					}
				}
			}
			if corrected != "" {
				fmt.Print(s, "->", corrected)
				fmt.Println()
			}
		}
	}
}
