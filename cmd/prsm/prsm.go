package main

import "flag"
import "fmt"
import "io"
import "log"
import "os"
import "rosalind/spectrum"

func main() {
	flag.Float64Var(&spectrum.ConvUnit, "convUnit", 1.0,
	                "Precision of spectral convolution")
	flag.Parse()

	prots, masses, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

    db := []*spectrum.Spectrum{}
	for _, pr := range prots {
		s, err := spectrum.FromProtein(pr)
		if err != nil {
			log.Fatal(err)
		}
		db = append(db, s)
	}

	spec, err := spectrum.New(masses)
	if err != nil {
		log.Fatal(err)
	}

	best := 0
	which := ""
	for _, dbSpec := range db {
		pr := dbSpec.SourceProtein
		conv := spec.Convolution(dbSpec)
		m, k := conv.Max()
		fmt.Println(pr)
		fmt.Println(k, m)
		if k >= best {
			best = k
			which = pr
		}
	}
	fmt.Println(best)
	fmt.Println(which)
}

func readInput(r io.Reader) ([]string, []spectrum.Mass, error) {
	var n int
	_, err := fmt.Fscanf(r, "%d\n", &n)
	if err != nil {
		return nil, nil, err
	}

	prots := []string{}
	for len(prots) < n {
		var s string
		_, err := fmt.Fscanf(r, "%s\n", &s)
		if err != nil {
			return nil, nil, err
		}
		prots = append(prots, s)
	}

    masses := []spectrum.Mass{}
	for {
		var s string
		_, err := fmt.Fscanf(r, "%s\n", &s)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}
		m, err := spectrum.ParseMass(s)
		if err != nil {
			return nil, nil, err
		}
		masses = append(masses, m)
	}

	return prots, masses, nil
}
