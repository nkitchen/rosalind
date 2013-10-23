package main

import "bufio"
import "fmt"
import "os"
import "rosalind/strings/suffix"
import "sort"

const minRepeatLen = 2

type SuffixTree suffix.Tree

type byLength []string

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()

	t := (*SuffixTree)(suffix.NewTree(s + "$"))

	repeats := []string{}
	t.forEachMaximalRepeat(1, 0, func(m string) {
		repeats = append(repeats, m)
	})
	sort.Sort(byLength(repeats))

	maximal := []string{}
Repeat:
	for _, r := range repeats {
		for _, m := range maximal {
			if m[len(m) - len(r):] == r {
				continue Repeat
			}
		}
		maximal = append(maximal, r)
	}

	for _, m := range maximal {
		fmt.Println(m)
	}
}

func (t *SuffixTree) forEachMaximalRepeat(node int, prefixLen int, f func (string)) (maximal bool) {
	for _, e := range t.Edges[node] {
		if e.LeafCount < 2 {
			maximal = true
			continue
		}

		i := e.Loc - prefixLen
		n := prefixLen + e.Len

		m := t.forEachMaximalRepeat(e.Head, n, f)
		if m && n >= minRepeatLen {
			f(t.String[i:i+n])
		}
	}
	return
}

func (a byLength) Len() int {
	return len(a)
}

func (a byLength) Less(i, j int) bool {
	return len(a[i]) > len(a[j])
}

func (a byLength) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
