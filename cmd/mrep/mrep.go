package main

import "bufio"
import "fmt"
import "os"
import "rosalind/strings"
import "rosalind/strings/suffix"

const minRepeatLen = 20

type SuffixTree suffix.Tree

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()

	tFwd := (*SuffixTree)(suffix.NewTree(s + "$"))

	fwdMaximal := map[string]bool{}
	tFwd.forEachMaximalRepeat(1, 0, func(m string) {
		fwdMaximal[m] = true
	})

	b := []byte(s)
	strings.ReverseBytes(b)
	r := string(b)
	tRev := (*SuffixTree)(suffix.NewTree(r + "$"))

	tRev.forEachMaximalRepeat(1, 0, func(m string) {
		b := []byte(m)
		strings.ReverseBytes(b)
		mr := string(b)
		if _, ok := fwdMaximal[mr]; ok {
			fmt.Println(mr)
		}
	})
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
