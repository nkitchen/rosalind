package graph

import "bufio"
import "io"
import "strconv"
import "strings"

type Error string
func (e Error) Error() string {
	return string(e)
}

// ReadDirectedAdjacencies reads an adjacency list for a directed graph
// from a reader.  It assumes that nothing follows the graph
// in the reader's data.
// It assumes that the nodes are numbered 1..n in the input,
// but the returned list has node numbers 0..n-1.
func ReadDirectedAdjacencies(r io.Reader) ([][]int, error) {
	br := bufio.NewReader(r)
	line, err := br.ReadString('\n')
	if err != nil {
		return nil, err
	}
	n, err := strconv.Atoi(strings.TrimSpace(line))

	g := make([][]int, n)
	for {
		line, err = br.ReadString('\n')
		if err == io.EOF {
			return g, nil
		} else if err != nil {
			return g, err
		}

		words := strings.Fields(line)
		if len(words) != 2 {
			return g, Error("Bad adjacency: " + strings.TrimSpace(line))
		}
		var edge [2]int
		for i, w := range words {
			v, err := strconv.Atoi(w)
			if err != nil {
				return g, Error("Bad node index: " + w)
			}
			edge[i] = v - 1
		}
		g[edge[0]] = append(g[edge[0]], edge[1])
	}
	return nil, Error("unreachable")
}

// Assumes that the graph is given as an adjacency list.
func ConnectedComponents(g [][]int) [][]int {
	rep := make([]int, len(g))
	for v := range g {
		rep[v] = v
	}

	var find func(v int) int
	find = func(v int) int {
		r := rep[v]
		if r != v {
			r = find(r)
			rep[v] = r
		}
		return r
	}

	for u, a := range g {
		for _, v := range a {
			rep[find(u)] = find(v)
		}
	}

	inv := make(map[int][]int)
	for v := range g {
		r := find(v)
		inv[r] = append(inv[r], v)
	}

	ccs := [][]int{}
	for _, cc := range inv {
		ccs = append(ccs, cc)
	}
	return ccs
}
