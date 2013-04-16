package strings

import "fmt"

func min(x int, ys ...int) int {
	m := x
	for _, y := range ys {
		if y < m {
			m = y
		}
	}
	return m
}

var DebugEditDistance = false

// Returns the edit distance between the bytes of two strings,
// which is the number of edit operations (insertions, deletions,
// or substitutions of single characters) needed to transform
// one string into the other.
func EditDistance(s, t string) int {
	m := len(s)
	n := len(t)

	a := make([][]int, m)
	for i := range a {
		a[i] = make([]int, n)
	}

	if s[0] != t[0] {
		a[0][0] = 1
	}

	for i := 1; i < m; i++ {
		if s[i] == t[0] {
			a[i][0] = i
		} else {
			a[i][0] = a[i - 1][0] + 1
		}
	}
	for j := 1; j < n; j++ {
		if s[0] == t[j] {
			a[0][j] = j
		} else {
			a[0][j] = a[0][j - 1] + 1
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			var diag int
			if s[i] == t[j] {
				diag = 0
			} else {
				diag = 1
			}
			d := min(a[i - 1][j] + 1, a[i][j - 1] + 1, a[i - 1][j - 1] + diag)
			a[i][j] = d
		}
	}

	if DebugEditDistance {
		fmt.Print("  ")
		for j := range t {
			fmt.Printf("%c  ", t[j])
		}
		for i := range a {
			fmt.Printf("%c ", s[i])
			for j := range a[i] {
				fmt.Printf("%2v ", a[i][j])
			}
			fmt.Println()
		}
	}

	return a[m - 1][n - 1]
}

var _ = fmt.Println
