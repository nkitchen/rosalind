package main

import "bufio"
import "fmt"
import "log"
import "math"
import "os"
import "strconv"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	dna := strings.TrimSpace(line)

	line, err = br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range strings.Fields(line) {
		gc, _ := strconv.ParseFloat(f, 64)
		lgc := math.Log10(gc / 2)
		lat := math.Log10((1 - gc) / 2)
		lp := float64(0)
		for i := 0; i < len(dna); i++ {
			switch dna[i] {
			case 'A', 'T':
			    lp += lat
			case 'G', 'C':
				lp += lgc
			}
		}
		fmt.Print(lp, " ")
	}
	fmt.Println()
} 
