package main

import "bufio"
import "fmt"
import "log"
import "os"
import "sort"
import "strconv"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	line = line[:len(line) - 1]

	n64, err := strconv.ParseInt(line, 0, 32)
	n := int(n64)

    p := make([]int, n)
	m := 1
	for i := 0; i < n; i++ {
		p[i] = i + 1
		m *= i + 1
	}

    fmt.Println(m)
	for {
		for i, x := range p {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(x)
		}
		fmt.Println()
		ok := lexicoNext(p)
		if !ok {
			break
		}
	}
} 

// lexicoNext reorders s to the next permutation in lexicographic order.
// Returns false if s is already the final permutation; otherwise true.
func lexicoNext(s []int) bool {
	var i int
	for i = len(s) - 2; i >= 0; i-- {
		if s[i] < s[i + 1] {
			break
		}
	}
	if i < 0 {
		return false
	}

    minGreater := i + 1
	for j := minGreater + 1; j < len(s); j++ {
		if s[j] > s[i] && s[j] < s[minGreater] {
			minGreater = j
		}
	}

	s[i], s[minGreater] = s[minGreater], s[i]
	sort.IntSlice(s[i+1:]).Sort()
	return true
}
