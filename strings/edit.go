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

const Gap = -1

type multAlignEntry struct {
	score int32
	// Points to the entry from which the score is derived:
	// If bit k is not set (decrBits & (1 << k) == 0), the kth index
	// of the previous entry is the same as the kth index of this
	// entry; if bit k is set, the kth index is one less.
	decrBits int32
}

type multAligner struct {
	strings []string
	memo map[string]multAlignEntry
	scoreFunc func(int16, int16) int32
}

func (aligner multAligner) Get(indices []int32) multAlignEntry {
	key := string(indices)
	if e, ok := aligner.memo[key]; ok {
		return e
	}

	origin := true
	for _, i := range indices {
		if i > 0 {
			origin = false
			break
		}
	}
	if origin {
		e := multAlignEntry{0, 0}
		aligner.memo[key] = e
		return e
	}

	bestScore := int32(0)
	bestDecrBits := int32(0)
	maxDecrBits := int64(1) << uint(len(indices)) - 1
DecrLoop:
	for decrBits := int32(1); int64(decrBits) <= maxDecrBits; decrBits++ {
		prevIndices := make([]int32, len(indices))
		copy(prevIndices, indices)
		for i := range indices {
			if decrBits & (1 << uint(i)) != 0 {
				if indices[i] == 0 {
					continue DecrLoop
				}
				prevIndices[i]--
			}
		}

		score := aligner.Get(prevIndices).score
		for i := range aligner.strings {
			var c, d int16
			if decrBits & (1 << uint(i)) == 0 {
				c = Gap
			} else {
				c = int16(aligner.strings[i][indices[i] - 1])
			}
			for j := i + 1; j < len(aligner.strings); j++ {
				if decrBits & (1 << uint(j)) == 0 {
					d = Gap
				} else {
					d = int16(aligner.strings[j][indices[j] - 1])
				}

				score += aligner.scoreFunc(c, d)
			}
		}
		if bestDecrBits == 0 || score > bestScore {
			bestScore = score
			bestDecrBits = decrBits
		}
	}

	e := multAlignEntry{bestScore, bestDecrBits}
	aligner.memo[key] = e
	return e
}

// Returns augmented strings for the elements of s so as to maximize the sum
// of scores over pairs of augmented strings.
// The arguments to scoreFunc can be bytes from two strings or Gap.
// gapSym is the symbol inserted into gaps in the augmented strings.
func MultipleAlignment(s []string, scoreFunc func(int16, int16) int32,
                       gapSym byte) (int32, []string) {
	if len(s) > 32 {
		panic("MultipleAlignment cannot operate on more than 32 strings.")
	}

	aligner := multAligner{s, make(map[string]multAlignEntry), scoreFunc}
	indices := make([]int32, len(s))
	for i := range s {
		indices[i] = int32(len(s[i]))
	}

	entry := aligner.Get(indices)
	score := entry.score

	b := make([][]byte, len(s))
	for entry.decrBits != 0 {
		for i := range s {
			if entry.decrBits & (1 << uint(i)) == 0 {
				b[i] = append(b[i], gapSym)
			} else {
				b[i] = append(b[i], s[i][indices[i] - 1])
				indices[i]--
			}
		}
		entry = aligner.Get(indices)
	}

	a := []string{}
	for i := range b {
		reverseBytes(b[i])
		a = append(a, string(b[i]))
	}

	return score, a
}

var _ = fmt.Println
