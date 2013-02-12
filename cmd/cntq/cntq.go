package main

import "bufio"
import "fmt"
import "os"
import "rosalind/tree"

const M = 1000000

var leaves = map[*tree.Node]int64{}

func main() {
	br := bufio.NewReader(os.Stdin)

	_, _ = br.ReadString('\n')

    t, _ := tree.ReadNewick(br)

	countLeaves(t)
	e := tree.Edge{t, 0}
	_, q := countPQ(e, 0, 0)
	fmt.Println(q)
}

func countLeaves(n *tree.Node) {
	if len(n.Children) == 0 {
		leaves[n] = 1
		return
	}

	k := int64(0)
	for _, e := range n.Children {
		countLeaves(e.Node)
		k += leaves[e.Node]
	}
	leaves[n] = k
}

var foo = 0
func countPQ(e tree.Edge, priorLeaves, priorPairs int64) (pairs, quartets int64) {
	bar := foo
	debugf("countPQ:%v priorPairs=%v\n", bar, priorPairs)
	foo++
	if debug {
		tree.Print(e.Node)
	}

	if leaves[e.Node] == 1 {
		return 0, 0
	}

	// pairs return value: their deepest common ancestor is the node below this
	// edge

	// pairs that include leaves in a particular subtree
	pairsWith := make([]int64, len(e.Children))
	for i := 0; i < len(e.Children) - 1; i++ {
		for j := i + 1; j < len(e.Children); j++ {
			p := leaves[e.Children[i].Node] * leaves[e.Children[j].Node]
			pairs += p
			pairsWith[i] += p
			pairsWith[j] += p
		}
	}
	debugf("pairsWith=%v\n", pairsWith)

	quartets += priorPairs * pairs

	for i, f := range e.Children {
		debugf(":%v[child %v]: pairs=%v quartets=%v\n", bar, i, pairs, quartets)
	    pp := priorPairs + priorLeaves * (leaves[e.Node] - leaves[f.Node])
		pl := priorLeaves + leaves[e.Node] - leaves[f.Node]
		p, q := countPQ(f, pl, pp + pairs - pairsWith[i])
		pairs += p
		quartets += q
	}

	debugf("ret:%v pairs=%v quartets=%v\n", bar, pairs, quartets)
	return
}

var debug = false
func debugf(fmt string, args ...interface {}) {
	if debug {
		debugf(fmt, args...)
	}
}
