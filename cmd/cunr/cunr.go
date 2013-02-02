package main

import "bufio"
import "fmt"
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
		_, b := binom(n, k)
		t := (b * rootedTrees(k)) % M
		//fmt.Println("k", k, "t", t)
		t = (t * rootedTrees(n - k)) % M
		s = (s + t) % M		
		//fmt.Println("k", k, "t", t, "s", s)
	}

	if k * 2 == n {
		r := rootedTrees(k)
		bq, br := binom(n, k)
		if br % 2 != 0 {
			panic(n)
		}
		t := br / 2 * r % M
		t = (t * r) % M
		s = (s + t) % M
		
		if bq % 2 == 1 {
			t = bq * M / 2 % M
			t = (t * r) % M
			t = (t * r) % M
			s = (s + t) % M
		}
		//fmt.Println("k", k, "r", r, "b", b, "t", t, "s", s)
	}

	treesMemo[n] = s
	return s
}

var binomMemo = map[[2]int64][2]int64{}
// Returns q, r such that the binomial coefficient of n and k is
// Q * M + r for some Q and Q % M == q.
func binom(n, k int64) (int64, int64) {
	s := [2]int64{n, k}
	m, ok := binomMemo[s]
	if ok {
		return m[0], m[1]
	}
	
	if k == 1 {
		return 0, n % M
	}
	if k == 0 || k == n {
		return 0, 1
	}

	aq, ar := binom(n - 1, k - 1)
	bq, br := binom(n - 1, k)
	q := (aq + bq + (ar + br) / M) % M
	r := (ar + br) % M
	binomMemo[s] = [2]int64{q, r}
	return q, r
}
