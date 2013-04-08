package main

import "bufio"
import "flag"
import "fmt"
import "log"
import "os"
import "rosalind/phylo"
import "rosalind/tree"
import "runtime/pprof"
import "sort"
import "strings"

func main() {
	verify := flag.Bool("verify", false, "Check the distance for correctness")
	memprof := flag.String("memprof", "", "Write memory profile to this file")
	flag.Parse()

	br := bufio.NewReader(os.Stdin)
	
	line, _ := br.ReadString('\n')
	taxa := strings.Fields(line)
	taxaInv := map[string]int{}
	for i, taxon := range taxa {
		taxaInv[taxon] = i
	}

	line, _ = br.ReadString('\n')
	t1, err := tree.ReadNewick(strings.NewReader(line))
	if err != nil {
		panic(err)
	}
	line, _ = br.ReadString('\n')
	t2, err := tree.ReadNewick(strings.NewReader(line))
	if err != nil {
		panic(err)
	}

	d := phylo.QuartetDistance(t1, t2, len(taxa))
	fmt.Println(d)

	if *verify {
		d2 := distanceFromSplits(t1, t2, taxaInv)
		fmt.Println("From splits:", d2)
		if d == d2 {
			fmt.Println("ok")
		} else {
			fmt.Println("mismatch")
			os.Exit(1)
		}
	}

	if *memprof != "" {
		f, err := os.Create(*memprof)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.Lookup("heap").WriteTo(f, 0)
	}
} 

func distanceFromSplits(t1, t2 *tree.Node, taxa map[string]int) int {
	s1 := phylo.Splits(t1, taxa)
	s2 := phylo.Splits(t2, taxa)

	q1 := quartetsFromSplits(s1)
	q2 := quartetsFromSplits(s2)

	shared := 0
	for q := range q1 {
		if q2[q] {
			shared++
		}
	}

	return len(q1) + len(q2) - 2 * shared
}

func quartetSlice(m map[phylo.Quartet]bool) []phylo.Quartet {
	s := make([]phylo.Quartet, 0, len(m))
	for q := range m {
		s = append(s, q)
	}
	sort.Sort(phylo.QuartetSlice(s))
	return s
}

func quartetsFromSplits(splits []phylo.CharArray) map[phylo.Quartet]bool {
	quartets := map[phylo.Quartet]bool{}

	for _, s := range splits {
		for a := 0; a < s.Len() - 1; a++ {
			if s.At(a) != 0 {
				continue
			}
			for b := a + 1; b < s.Len(); b++ {
				if s.At(b) != 0 {
					continue
				}
				p1 := phylo.NewPair(a, b)

				for c := 0; c < s.Len() - 1; c++ {
					if s.At(c) == 0 {
						continue
					}
					for d := c + 1; d < s.Len(); d++ {
						if s.At(d) == 0 {
							continue
						}
						p2 := phylo.NewPair(c, d)
						quartets[phylo.NewQuartet(p1, p2)] = true
					}
				}
			}
		}
	}

	return quartets
}
