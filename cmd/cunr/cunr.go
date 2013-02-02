package main

import "bufio"
import "fmt"
import "math/big"
import "os"
import "strconv"
import "strings"

const M = 1000000

func main() {
	br := bufio.NewReader(os.Stdin)

	line, _ := br.ReadString('\n')
	n, _ := strconv.ParseInt(strings.TrimSpace(line), 0, 64)

	//for m := int64(4); m <= n; m++ {
	//	fmt.Println(m, rootedTrees(m - 1))
	//}
	fmt.Println(rootedTrees(n - 1))
}

var bigM = big.NewInt(M)
var treesMemo = map[int64]int64{}
func rootedTrees(n int64) int64 {
	if n == 1 {
		return 1
	}

	s := treesMemo[n]
	if s != 0 {
		return s
	}

	s = int64(0)
	var k int64
	for k = 1; 2 * k < n; k++ {
		t := (binom(n, k) * rootedTrees(k)) % M
		//fmt.Println("k", k, "t", t)
		t = (t * rootedTrees(n - k)) % M
		s = (s + t) % M		
		//fmt.Println("k", k, "t", t, "s", s)
	}

	if k * 2 == n {
		r := rootedTrees(k)
		b := big.NewInt(0)
		b.Binomial(n, k)
		var quo, rem big.Int
		quo.QuoRem(b, bigM, &rem)
		if rem.Bit(0) != 0 {
			panic(n)
		}
		t := rem.Int64() / 2 * r % M
		t = (t * r) % M
		s = (s + t) % M
		if quo.Bit(0) == 1 {
			quo.Rem(&quo, bigM)
			t = quo.Int64() * M / 2 % M
			t = (t * r) % M
			t = (t * r) % M
			s = (s + t) % M
		}
		//fmt.Println("k", k, "r", r, "b", b, "t", t, "s", s)
	}

	treesMemo[n] = s
	return s
}

var binomMemo = map[[2]int64]int64{}
func binom(n, k int64) int64 {
	s := [2]int64{n, k}
	r := binomMemo[s]
	if r != 0 {
		return r
	}
	
	if k == 1 {
		return n % M
	}
	if k == 0 || k == n {
		return 1
	}

	a := binom(n - 1, k - 1)
	b := binom(n - 1, k)
	r = (a + b) % M
	binomMemo[s] = r
	return r
}
