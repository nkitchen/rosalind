package main

import "bufio"
import "fmt"
import "log"
import "os"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	s := line[:len(line) - 1]

	line, err = br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	t := line[:len(line) - 1]

	fmt.Println(HammingDist(s, t))
}

func HammingDist(s, t string) int {
	if len(s) != len(t) {
		log.Fatalf("Mismatched string lengths: %v vs. %v\n", len(s), len(t))
	}
	d := 0
	for i := range(s) {
		if s[i] != t[i] {
			d++
		}
	}
	return d
}

//
//type DnaString struct {
//	Id string
//	Bases string
//}
//
//func ReadAllFasta(r io.Reader) []DnaString {
//	var dna []DnaString
//	br := bufio.NewReader(r)
//	for {
//		line, err := br.ReadString('\n')
//		if (err != nil && err != io.EOF) {
//			return dna
//		}
//
//		if (len(line) == 0) {
//			return dna
//		}
//
//		n := len(line) - 1
//		if line[n] == '\n' {
//			line = line[:n]
//		}
//
//		if (line[0] == '>') {
//			dna = append(dna, DnaString{Id: line[1:]})
//		} else {
//			p := &dna[len(dna)-1].Bases
//			*p += line
//		}
//
//		if err == io.EOF {
//			return dna
//		}
//	}
//	return nil
//}
//
//func gcContent(dna DnaString) float64 {
//	gc := float64(0)
//	for _, b := range dna.Bases {
//		switch b {
//		case 'C', 'G': gc += 1
//		}
//	}
//	n := float64(len(dna.Bases))
//	return gc / n
//}
//
//func main() {
//	var bestId string
//	var bestGc float64
//	for _, dna := range ReadAllFasta(os.Stdin) {
//		gc := gcContent(dna)
//		if bestId == "" || gc > bestGc {
//			bestId = dna.Id
//			bestGc = gc
//		}
//	}
//	fmt.Println(bestId)
//	fmt.Printf("%.5f%%\n", bestGc * 100)
//}
