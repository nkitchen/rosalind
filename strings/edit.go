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
	a := editMatrix(s, t)
	return a[len(s) - 1][len(t) - 1]
}

func editMatrix(s, t string) [][]int {
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
		printMatrix(s, t, a)
	}

	return a
}

func printMatrix(s, t string, a [][]int) {
	fmt.Print("  ")
	if len(a[0]) == len(t) + 1 {
		fmt.Print("    ")
	}
	for j := range t {
		fmt.Printf("  %c ", t[j])
	}
	fmt.Println()
	for i := range a {
		if len(a) == len(s) + 1 {
			if i == 0 {
				fmt.Print("  ")
			} else {
				fmt.Printf("%c ", s[i - 1])
			}
		} else {
			fmt.Printf("%c ", s[i])
		}

		for j := range a[i] {
			fmt.Printf("%3v ", a[i][j])
		}
		fmt.Println()
	}
}

// Returns supersequences of s and t obtained by inserting the gap symbol
// into them such that a maximum number of bytes of s and t are at the
// same positions.
func Alignment(s, t string, gapSym byte) (string, string) {
	a := editMatrix(s, t)

	// We construct the supersequences in reverse to avoid deep recursion
	// on the matrix.
	sa := make([]byte, 0, len(s))
	ta := make([]byte, 0, len(t))

	i := len(s) - 1
	j := len(t) - 1
	for i >= 0 && j >= 0 {
		switch {
		case i > 0 && j > 0 && s[i] == t[j] && a[i][j] == a[i - 1][j - 1],
		     i == 0 && j == 0:
			sa = append(sa, s[i])
			ta = append(ta, t[j])
			i--
			j--
		case j > 0 && a[i][j] == a[i][j - 1] + 1:
			sa = append(sa, gapSym)
			ta = append(ta, t[j])
			j--
		case i > 0 && a[i][j] == a[i - 1][j] + 1:
			sa = append(sa, s[i])
			ta = append(ta, gapSym)
			i--
		case i > 0 && j > 0 && s[i] != t[j] && a[i][j] == a[i - 1][j - 1] + 1:
			sa = append(sa, s[i])
			ta = append(ta, t[j])
			i--
			j--
		default:
			panic("Unexpected case")
		} 
	}

	reverseBytes(sa)
	reverseBytes(ta)
	return string(sa), string(ta)
}

func reverseBytes(a []byte) {
	i := 0
	j := len(a) - 1
	for i < j {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}
}

// Returns the maximum alignment score.
// scoringMatrix[i][j] is the value of replacing byte i in s with byte j in t.
// gapPenalty is the deduction from the score for each unmatched byte.
func MaxAlignmentScore(s, t string, scoringMatrix [][]int, gapPenalty int) int {
	m := len(s)
	n := len(t)

	a := make([][]int, m + 1)
	for i := 0; i <= m; i++ {
		a[i] = make([]int, n + 1)
	}

	for i := range s {
		a[i + 1][0] = (i + 1) * gapPenalty
	}
	for j := range t {
		a[0][j + 1] = (j + 1) * gapPenalty
	}

	for i := range s {
		for j := range t {
			c := int(s[i])
			d := int(t[j])
			s1 := a[i][j] + scoringMatrix[c][d]
			s2 := a[i + 1][j] + gapPenalty
			s3 := a[i][j + 1] + gapPenalty
			score := s1
			if s2 > score {
				score = s2
			}
			if s3 > score {
				score = s3
			}
			a[i + 1][j + 1] = score
		}
	}

	if DebugEditDistance {
		printMatrix(s, t, a)
	}

	return a[m][n]
}

var _ = fmt.Println
