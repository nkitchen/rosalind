package main

import "fmt"
import "math/big"
import "os"
import "rosalind/gene"

var Complement = map[byte]byte{
	'A': 'U',
	'C': 'G',
	'G': 'C',
	'U': 'A',
}

func main() {
	s, _ := gene.ReadFasta(os.Stdin)

	var matchings func(int, int) *big.Int
	memo := make(map[[2]int]*big.Int)
	matchings = func(begin, end int) *big.Int {
		switch {
		case begin == end:
			return big.NewInt(1)
		case begin > end:
			return big.NewInt(0)
		}

		key := [2]int{begin, end}
		if n, ok := memo[key]; ok {
			return n
		}

		n := &big.Int{}
		n.Set(matchings(begin + 1, end))
		compl := Complement[s.Data[begin]]
		for i := begin + 1; i < end; i++ {
			if s.Data[i] == compl {
				var t big.Int
				t.Mul(matchings(begin + 1, i), matchings(i + 1, end))
				n.Add(n, &t)
			}
		}
		memo[key] = n
		return n
	}

    n := matchings(0, len(s.Data))
	//n1m := big.NewInt(1000000)
	//n1m.Rem(n, n1m)
	//fmt.Println(n1m)
	fmt.Println(n)
}
