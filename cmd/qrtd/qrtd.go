package main

import "bufio"
import "flag"
import "fmt"
import "os"
import "rosalind/phylo"
import "rosalind/tree"
import "sort"
import "strings"

func main() {
	verify := flag.Bool("verify", false, "Check the distance for correctness")
	flag.Parse()

	br := bufio.NewReader(os.Stdin)
	
	line, _ := br.ReadString('\n')
	taxa := strings.Fields(line)
	taxaInv := map[string]int{}
	for i, taxon := range taxa {
		taxaInv[taxon] = i
	}

	line, _ = br.ReadString('\n')
	t1, _ := tree.ReadNewick(strings.NewReader(line))
	line, _ = br.ReadString('\n')
	t2, _ := tree.ReadNewick(strings.NewReader(line))

	d := phylo.QuartetDistance(t1, t2, taxaInv)
	fmt.Println(d)

	if *verify {
		d2 := distanceFromSplits(t1, t2, taxaInv)
		fmt.Println("From splits:", d2)
		if d == d2 {
			fmt.Println("ok")
		} else {
			fmt.Println("mismatch")
		}
	}
} 

func distanceFromSplits(t1, t2 *tree.Node, taxa map[string]int) int {
	s1 := phylo.Splits(t1, taxa)
	s2 := phylo.Splits(t2, taxa)

	fmt.Println("splits1:")
	for _, a := range s1 {
		fmt.Println(a)
	}
	fmt.Println("splits2:")
	for _, a := range s2 {
		fmt.Println(a)
	}

	q1 := quartetsFromSplits(s1)
	q2 := quartetsFromSplits(s2)

	fmt.Println("quartets1:", quartetSlice(q1))
	fmt.Println("quartets2:", quartetSlice(q2))

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
		for a := 0; a < len(s) - 1; a++ {
			if s[a] != 0 {
				continue
			}
			for b := a + 1; b < len(s); b++ {
				if s[b] != 0 {
					continue
				}
				p1 := phylo.NewPair(a, b)

				for c := 0; c < len(s) - 1; c++ {
					if s[c] == 0 {
						continue
					}
					for d := c + 1; d < len(s); d++ {
						if s[d] == 0 {
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
