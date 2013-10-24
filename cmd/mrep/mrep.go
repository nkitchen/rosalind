package main

import "bufio"
import "flag"
import "fmt"
import "os"
import "rosalind/strings/suffix"

var minRepeatLen int

type SuffixTree suffix.Tree

var edgeLocs map[*suffix.Edge][]int
var repeats map[string]map[int]bool

func main() {
	flag.IntVar(&minRepeatLen, "len", 20, "Minimum repeat length")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()

	t := (*SuffixTree)(suffix.NewTree(s + "$"))

	edgeLocs = map[*suffix.Edge][]int{}
	t.findEdgeLocs(1)

	repeats = map[string]map[int]bool{}
	t.findRepeats(1, 0)

	for rep, locs := range repeats {
		pre := map[rune]bool{}
		post := map[rune]bool{}
		for i := range locs {
			if i == 0 {
				pre['^'] = true
			} else {
				c := rune(t.String[i - 1])
				pre[c] = true
			}

			if i + len(rep) < len(s) {
				c := rune(t.String[i + len(rep)])
				post[c] = true
			} else {
				post['$'] = true
			}
		}
		if len(pre) != 1 && len(post) != 1 {
			fmt.Println(rep)
		}
	}
}

func (t *SuffixTree) findEdgeLocs(node int) []int {
	edges := t.Edges[node]
	if len(edges) == 0 {
		return []int{len(t.String)}
	}

	locs := []int{}
	for _, e := range edges {
		for _, i := range t.findEdgeLocs(e.Head) {
			edgeLocs[e] = append(edgeLocs[e], i - e.Len)
			locs = append(locs, i - e.Len)
		}
	}
	return locs
}

func (t *SuffixTree) findRepeats(node int, prefixLen int) {
	for _, e := range t.Edges[node] {
		if len(edgeLocs[e]) < 2 {
			continue
		}

		n := prefixLen + e.Len
		t.findRepeats(e.Head, n)
		if n < minRepeatLen {
			continue
		}

		for _, i := range edgeLocs[e] {
			repLoc := i - prefixLen
			rep := t.String[repLoc:repLoc + n]
			if repeats[rep] == nil {
				repeats[rep] = map[int]bool{}
			}
			repeats[rep][repLoc] = true
		}
	}
}
