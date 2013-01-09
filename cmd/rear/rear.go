package main

import "bufio"
import "container/heap"
import "fmt"
import "io"
import "log"
import "os"
import "strings"
import "strconv"

func main() {
	br := bufio.NewReader(os.Stdin)

Input:
	for {
		p, err := readPerm(br)
		switch err {
			case io.EOF: break Input
			case nil:
			default: log.Fatal(err)
		}
		q, err := readPerm(br)
		if err != nil {
			log.Fatal(err)
		}
			
		fmt.Println(reversalDist(p, q))
	}
}

func readPerm(br *bufio.Reader) ([]int, error) {
	line, err := br.ReadString('\n')
	for err == nil && strings.TrimSpace(line) == "" {
		line, err = br.ReadString('\n')
	}

    if err != nil {
		return nil, err
	}

	words := strings.Fields(strings.TrimSpace(line))
	p := make([]int, 0, len(words))
	for _, w := range words {
		n, err := strconv.ParseInt(w, 0, 32)
		if err != nil {
			return p, err
		}
		p = append(p, int(n))
	}
	return p, nil
}

type node struct {
	// The permutation
	perm []int
	// The indices where consecutive elements are not consecutive values
	skips []int
	// The number of reversals applied so far
	reversals int
}

type queue []*node

func reversalDist(p, q []int) int {
	if len(p) != len(q) {
		log.Fatalf("Mismatched lengths: %v != %v", len(p), len(q))
		return -1
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
	n := node{p, skips(p), 0}
	heap.Push(&prioQ, &n)
	for len(prioQ) > 0 {
		n := heap.Pop(&prioQ).(*node)
		fmt.Println(*n)
		if len(n.skips) == 0 {
			return n.reversals
		}

		inv := map[int]int{-1: -1, len(n.perm): len(n.perm)}
		for i, x := range n.perm {
			inv[x] = i
		}
		for i, a := range n.skips {
			j := inv[a - 1]
			if a > 0 && j < i {
				q := reverse(n.perm, j, i - 1)
				nx := node{q, skips(q), n.reversals + 1}
				heap.Push(&prioQ, &nx)
			}

			var b int
			if i > 0 {
				b = n.perm[i - 1]
			} else {
				b = -1
			}
			j = inv[b + 1]
			if b < len(n.perm) - 1 && j > i {
				q := reverse(n.perm, i, j)
				nx := node{q, skips(q), n.reversals + 1}
				heap.Push(&prioQ, &nx)
			}
		}
	}
	return -1
}

func reverse(p []int, i, j int) []int {
	q := make([]int, len(p))
	copy(q, p)
	for i < j {
		q[i], q[j] = q[j], q[i]
		i++
		j--
	}
	return q
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
	return 2 * u.reversals + len(u.skips) < 2 * v.reversals + len(v.skips)
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
	*q = (*q)[:n]
	return x
}
