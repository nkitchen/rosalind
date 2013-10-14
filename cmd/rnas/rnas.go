package main

import "bufio"
import "fmt"
import "math/big"
import "os"

type Pair [2]byte

var Bonds = map[[2]byte]bool {
	Pair{'A', 'U'}: true,
	Pair{'C', 'G'}: true,
	Pair{'G', 'C'}: true,
	Pair{'U', 'A'}: true,
	Pair{'G', 'U'}: true,
	Pair{'U', 'G'}: true,
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()

	var matchings func(int, int) *big.Int
	memo := make(map[[2]int]*big.Int)
	matchings = func(begin, end int) (r *big.Int) {
		switch {
		case begin > end:
			return big.NewInt(0)
		case begin + 4 > end:
			return big.NewInt(1)
		}

		key := [2]int{begin, end}
		if n, ok := memo[key]; ok {
			return n
		}

		n := &big.Int{}
		n.Set(matchings(begin + 1, end))
		for i := begin + 4; i < end; i++ {
			if Bonds[Pair{s[begin], s[i]}] {
				var t big.Int
				t.Mul(matchings(begin + 1, i), matchings(i + 1, end))
				n.Add(n, &t)
			}
		}
		memo[key] = n
		return n
	}

    n := matchings(0, len(s))
	fmt.Println(n)
}
