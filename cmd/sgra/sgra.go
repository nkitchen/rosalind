package main

import "bufio"
import "fmt"
import "log"
import "os"
import "rosalind/spectrum"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	masses := []spectrum.Mass{}
	for scanner.Scan() {
		w := scanner.Text()
		m, err := spectrum.ParseMass(w)
		if err != nil {
			log.Fatal(err)
		}
		masses = append(masses, m)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	spec := spectrum.New(masses)
	pr := spec.FindProtein()
	fmt.Println(pr)
}
