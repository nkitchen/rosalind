package main

import "bufio"
import "fmt"
import "os"

var s string

type Edge struct {
	Head int
	Loc int
	Len int
}

var edges map[int]map[byte]Edge
var nextNode = 1

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s = scanner.Text()

	edges = make(map[int]map[byte]Edge)
	edges[1] = make(map[byte]Edge)
	edges[1][s[0]] = Edge{2, 0, len(s)}
	nextNode = 3

	for i := 1; i < len(s); i++ {
		t := s[i:]
		tLoc := i
		node := 1
Prefix:
		for {
			e, ok := edges[node][t[0]]
			if !ok {
				edges[node][t[0]] = Edge{nextNode, tLoc, len(t)}
				nextNode++
				break
			}

			u := s[e.Loc:e.Loc + e.Len]
			if len(u) < len(t) && t[:len(u)] == u {
				t = t[len(u):]
				tLoc += len(u)
				node = e.Head
				continue
			}

			if t[0] != u[0] {
				panic("Logic error")
			}
			for j := 1; j < len(t); j++ {
				if t[j] != u[j] {
					v := nextNode
					nextNode++
					w := nextNode
					nextNode++

					edges[node][t[0]] = Edge{v, e.Loc, j}
					edges[v] = make(map[byte]Edge)
					edges[v][s[e.Loc + j]] = Edge{e.Head, e.Loc + j, e.Len - j}
					edges[v][t[j]] = Edge{w, tLoc + j, len(t) - j}
					break Prefix
				}
			}
			panic("Unexpected match of t")
		}
	}

	for _, a := range edges {
		for _, e := range a {
			fmt.Println(s[e.Loc:e.Loc + e.Len])
		}
	}
}

