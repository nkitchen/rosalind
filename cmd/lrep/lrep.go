package main

import "bufio"
import "fmt"
import "os"
import "strconv"
import "strings"

type Edge struct {
	Head string
	Loc int
	Len int
	Occurrences int
}

var edges map[string][]*Edge

var s string
var k int
var longestSub map[string]string

func main() {
	edges = make(map[string][]*Edge)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s = scanner.Text()

	scanner.Scan()
	k, _ = strconv.Atoi(scanner.Text())

	for scanner.Scan() {
		r := strings.NewReader(scanner.Text())
		var parent, child string
		var loc, len int
		_, err := fmt.Fscanf(r, "%s %s %d %d", &parent, &child, &loc, &len)
		loc -= 1
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		edges[parent] = append(edges[parent], &Edge{child, loc, len, 0})
	}

    countOccurrences("node1")

	longestSub = make(map[string]string)
	findLongest("node1")
	for tail, sub := range longestSub {
		fmt.Println(tail, len(sub))
	}

	t := ""
	for _, sub := range longestSub {
		if len(sub) > len(t) {
			t = sub
		}
	}
	fmt.Println(t)
    //fmt.Println("digraph lrep {")
	//for tail, out := range edges {
	//	for _, edge := range out {
	//		i := edge.Loc
	//		n := edge.Len
	//		fmt.Printf("  %s -> %s [label=\"%d:+%d / %d\"];\n", tail, edge.Head, i, n, edge.Occurrences)
	//	}
	//}
	//fmt.Println("}")
}

func countOccurrences(tail string) int {
	if len(edges[tail]) == 0 {
		return 1
	}
	n := 0
	for _, edge := range edges[tail] {
		m := countOccurrences(edge.Head)
		edge.Occurrences = m
		n += m
	}
	return n
}

func findLongest(tail string) {
	longestBelowTail := ""
	for _, edge := range edges[tail] {
		if edge.Occurrences < k {
			continue
		}
		edgeSub := s[edge.Loc:edge.Loc + edge.Len]
		if len(edgeSub) > len(longestBelowTail) {
			longestBelowTail = edgeSub
		}
		findLongest(edge.Head)
		if t, ok := longestSub[edge.Head]; ok {
			if edge.Len + len(t) > len(longestBelowTail) {
				longestBelowTail = edgeSub + t
			}
		}
	}
	if longestBelowTail != "" {
		longestSub[tail] = longestBelowTail
	}
}
