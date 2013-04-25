package main

import "bufio"
import "fmt"
import "log"
import "math"
import "os"
import "strings"
import "rosalind/gene"
//import rstrings "rosalind/strings"
import "rosalind/tree"

var n int
var dna = map[string]string{}
// For each position,
// for each choice of symbol at this node,
// the indices of the symbols for the children that gave the best distance
var whereBest = map[*tree.Node][][5][2]int32{}

func main() {
	br := bufio.NewReader(os.Stdin)

	line, _ := br.ReadString('\n')
	t, err := tree.ReadNewick(strings.NewReader(line))
	if err != nil {
		log.Fatal(err)
	}

	fasta, err := gene.ReadAllFasta(br)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range fasta {
		dna[f.Description] = f.Data
	}
	n = len(fasta[0].Data)

	d := bestDistances(t)
	where := make([]int32, n)
	total := float32(0)
	for i := 0; i < n; i++ {
		bestDist := d[i][0]
		where[i] = 0
		for j := int32(1); j < 5; j++ {
			if d[i][j] < bestDist {
				bestDist = d[i][j]
				where[i] = j
			}
		}
		total += bestDist
	}

	fmt.Println(total)
	collectBestDna(t, where)
}

var alphabet = []byte("-ACGT")


// Returns the minimum Hamming distance for the subtree given each possible
// choice of symbol at each position.
func bestDistances(t *tree.Node) [][5]float32 {
	if len(t.Children) == 0 {
		d := dna[t.Label]
		b := make([][5]float32, n)
		for i := 0; i < n; i++ {
			for j := 0; j < 5; j++ {
				if d[i] != alphabet[j] {
					b[i][j] = float32(math.Inf(1))
				}
			}
		}
		return b
	}

	childBest0 := bestDistances(t.Children[0].Node)
	childBest1 := bestDistances(t.Children[1].Node)
	best := make([][5]float32, n)
	whereBest[t] = make([][5][2]int32, n)
	for i := 0; i < n; i++ {
		tab := [5][5]float32{}
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				s := childBest0[i][j] + childBest1[i][k]
				if j != k {
					s++
				}
				tab[j][k] = s
			}
		}

		for j := int32(0); j < 5; j++ {
			b := tab[j][j]
			where := [2]int32{j, j}
			for k := int32(0); k < 5; k++ {
				if tab[j][k] < b {
					b = tab[j][k]
					where[0], where[1] = j, k
				}
				if tab[k][j] < b {
					b = tab[k][j]
					where[0], where[1] = k, j
				}
			}
			best[i][j] = b
			whereBest[t][i][j] = where
		}
	}
	return best
}

func collectBestDna(t *tree.Node, where []int32) {
	if _, ok := dna[t.Label]; ok {
		return
	}

	bestDna := make([]byte, n)
	for i := range where {
		bestDna[i] = alphabet[where[i]]
	}
	dna[t.Label] = string(bestDna)
	fmt.Printf(">%s\n", t.Label)
	fmt.Println(dna[t.Label])

	b := [2][]int32{}
	for j := 0; j < 2; j++ {
		b[j] = make([]int32, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < 2; j++ {
			b[j][i] = whereBest[t][i][where[i]][j]
		}
	}
	for j := 0; j < 2; j++ {
		collectBestDna(t.Children[j].Node, b[j])
	}
}

