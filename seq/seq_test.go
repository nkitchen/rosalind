package seq_test

import "fmt"
import "math/rand"
import "rosalind/seq"
import "testing"

func TestLongestIncSubseq(t *testing.T) {
	for k := 0; k < 100; k++ {
		fmt.Print(".")
		// Random permutation
		n := 10 + rand.Intn(11)
		p := make([]int, n)
		for i := range p {
			p[i] = i + 1
		}
		// Shuffle
		for i := 0; i < len(p) - 1; i++ {
			r := rand.Intn(len(p) - 1 - i)
			j := i + 1 + r
			p[i], p[j] = p[j], p[i]
		}

		fast := seq.LongestIncreasingSubseqInts(p)
		slow := seq.SlowLGIS(p)
		if len(fast) != len(slow) {
			t.Fatalf("LGIS(%v): expected %v, got %v", p, slow, fast)
		}
	}
}
