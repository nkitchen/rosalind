package main

import "bufio"
import "fmt"
import "log"
import "os"
import "rosalind/spectrum"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	spec := []spectrum.Mass{}

	for scanner.Scan() {
		t := scanner.Text()
		w, err := spectrum.ParseMass(t)
		if err != nil {
			log.Fatal(err)
		}

		spec = append(spec, w)
	}

	fmt.Println(ProteinFromSpectrum(spec))
}

func ProteinFromSpectrum(spec []spectrum.Mass) string {
	p := []byte{}
	for i := 1; i < len(spec); i++ {
		d := spec[i] - spec[i - 1]
		aa, ok := spectrum.ResidueByMass(d)
		if !ok {
			log.Fatalf("Cannot find residue for %v, %v: %v",
			           spec[i - 1], spec[i], d)
		}
		p = append(p, aa)
	}

	return string(p)
}
