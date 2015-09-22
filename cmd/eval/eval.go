package main

import "bufio"
import "fmt"
import "log"
import "math"
import "os"
import "strconv"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	words := []string{}
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

    n, err := strconv.ParseFloat(words[0], 64)
	if err != nil {
		log.Fatal(err)
	}

	s := words[1]

	a := []float64{}
	for _, w := range words[2:] {
		x, err := strconv.ParseFloat(w, 64)
		if err != nil {
			log.Fatal(err)
		}
		a = append(a, x)
	}

    m := map[rune]float64{}
	for _, c := range s {
		m[c] = m[c] + 1
	}
	nAT := m['A'] + m['T']
	nGC := m['C'] + m['G']

	l2 := math.Log(2)
    for _, pGC := range a {
		pAT := 1 - pGC
		lpGC := math.Log(pGC)
		lpAT := math.Log(pAT)
		lps := nAT * (lpAT - l2) + nGC * (lpGC - l2)
		k := n - float64(len(s)) + 1
		lk := math.Log(k)
		p := math.Exp(lk + lps)
		//fmt.Println(pGC, pAT, lpGC, lpAT, l2, lps, k)
		//fmt.Println(p)
		fmt.Printf("%.9f ", p)
	}
	fmt.Println()
}
