package suffix

type Edge struct {
	Head int
	// The index in the string where this substring starts
	Loc int
	// The length of the substring
	Len int
}

type Tree struct {
	String string
	Edges map[int]map[byte]*Edge
}

func NewTree(s string) *Tree {
	t := &Tree{s, map[int]map[byte]*Edge{}}
	t.Edges[1] = map[byte]*Edge{}
	t.Edges[1][s[0]] = &Edge{2, 0, len(s)}
	nextNode := 3

	for i := 1; i < len(s); i++ {
		suffix := s[i:]
		sufLoc := i
		node := 1
Prefix:
		for {
			e, ok := t.Edges[node][suffix[0]]
			if !ok {
				t.Edges[node][suffix[0]] = &Edge{nextNode, sufLoc, len(suffix)}
				nextNode++
				break
			}

			u := s[e.Loc:e.Loc + e.Len]
			if len(u) < len(suffix) && suffix[:len(u)] == u {
				suffix = suffix[len(u):]
				sufLoc += len(u)
				node = e.Head
				continue
			}

			if suffix[0] != u[0] {
				panic("Unexpected mismatch in first characters of suffix " +
				      "and edge substring")
			}
			for j := 1; j < len(suffix); j++ {
				if suffix[j] != u[j] {
					v := nextNode
					nextNode++
					w := nextNode
					nextNode++

					t.Edges[node][suffix[0]] = &Edge{v, e.Loc, j}
					t.Edges[v] = map[byte]*Edge{}
					t.Edges[v][s[e.Loc + j]] = &Edge{e.Head, e.Loc + j, e.Len - j}
					t.Edges[v][suffix[j]] = &Edge{w, sufLoc + j, len(suffix) - j}
					break Prefix
				}
			}
			panic("Unexpected match of suffix")
		}
	}

	return t
}
