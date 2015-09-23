package seq

import "container/heap"
import "fmt"
import "sort"

type Swapper interface {
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

// Reverses the order of elements in the subsequence data[i:j].
func ReverseSub(data Swapper, i, j int) {
	j--
	for i < j {
		data.Swap(i, j)
		i++
		j--
	}
}

func reversed(p []int, i, j int) []int {
	q := make([]int, len(p))
	copy(q, p)
	ReverseSub(sort.IntSlice(q), i, j)
	return q
}

type node struct {
	// The permutation
	perm []int
	// The indices where consecutive elements are not consecutive values
	skips []int
	// The reversals applied so far
	reversals [][2]int
}

func (n *node) Cost() int {
	return 2 * len(n.reversals) + len(n.skips)
}

type queue []*node

func keyStr(p []int) string {
	a := make([]rune, len(p))
	for i, x := range p {
		a[i] = rune(x)
	}
	return string(a)
}

// Returns a sequence of reversals that convert p into q.
// Each reversal is encoded as a pair (i, j) such that
// s[i:j] is the subsequence reversed.
func Reversals(p, q []int) ([][2]int, error) {
	// A* search
	// Heuristic: 2 * reversals so far + remaining discontinuities

	if len(p) != len(q) {
		err := fmt.Errorf("Mismatched lengths: %v != %v", len(p), len(q))
		return nil, err
	}

	// Remap the elements so that we can search for 0..n-1.
	tr := map[int]int{}
	for i, x := range p {
		tr[x] = i
	}
	p = make([]int, len(q))
	for i, x := range q {
		p[i] = tr[x]
	}

    var prioQ queue
	visited := map[string]bool{}
	n := node{p, skips(p), nil}
	heap.Push(&prioQ, &n)
	for len(prioQ) > 0 {
		n := heap.Pop(&prioQ).(*node)
		if len(n.skips) == 0 {
			r := n.reversals
			ReverseSub(revSlice(r), 0, len(r))
			return r, nil
		}

		rev := map[[2]int]bool{}
		expand := func (i, j int) {
			r := [2]int{i, j}
			if rev[r] {
				return
			}
			q := reversed(n.perm, i, j)
			s := keyStr(q)
			if visited[s] {
				return
			}
			visited[s] = true
			nx := node{q, skips(q), append(n.reversals, r)}
			heap.Push(&prioQ, &nx)
			rev[r] = true
		}

		for _, i := range n.skips {
			// Skip from [i-1]: a to [i]: b
			var a, b int
			if i == 0 {
				a = -1
			} else {
				a = n.perm[i - 1]
			}
			if i < len(n.perm) {
				b = n.perm[i]
			} else {
				b = len(n.perm)
			}
            
			// Can a + 1 be reversed down to [i]?
			for j := i + 1; j < len(n.perm); j++ {
				if n.perm[j] == a + 1 {
					expand(i, j + 1)
					break
				}
			}

			// Can b - 1 be reversed up to [i-1]?
			for j := 0; j < i - 1; j++ {
				if n.perm[j] == b - 1 {
					expand(j, i)
					break
				}
			}
		}
	}
	return nil, fmt.Errorf("No solution found")
}

func skips(p []int) []int {
	s := []int{}
	if p[0] != 0 {
		s = append(s, 0)
	}
	for i := 1; i < len(p); i++ {
		d := p[i] - p[i-1]
		if d * d != 1 {
			s = append(s, i)
		}
	}
	if p[len(p)-1] != len(p) - 1 {
		s = append(s, len(p))
	}
	return s
}

func (q *queue) Len() int {
	return len(*q)
}

func (q *queue) Less(i, j int) bool {
	u := (*q)[i]
	v := (*q)[j]
	cu := u.Cost()
	cv := v.Cost()
	if cu < cv {
		return true
	}
	if cu == cv {
		return len(u.reversals) > len(v.reversals)
	}
	return false
}

func (pq *queue) Swap(i, j int) {
	q := *pq
	q[i], q[j] = q[j], q[i]
}

func (q *queue) Push(x interface{}) {
	*q = append(*q, x.(*node))
}

func (q *queue) Pop() interface{} {
	n := len(*q) - 1
	x := (*q)[n]
	(*q)[n] = nil
	*q = (*q)[:n]
	return x
}

type revSlice [][2]int

func (r revSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
