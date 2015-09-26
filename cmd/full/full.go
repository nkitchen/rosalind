package main

import "bufio"
import "fmt"
import "log"
import "os"
import "rosalind/spectrum"

func main() {
	masses := []spectrum.Mass{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		m, err := spectrum.ParseMass(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		masses = append(masses, m)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	spec, err := spectrum.New(masses[1:])
	if err != nil {
		log.Fatal(err)
	}
	p, ok := spec.Protein()
	if ok {
		fmt.Println(p)
	} else {
		log.Fatal("Protein not found")
	}
}
