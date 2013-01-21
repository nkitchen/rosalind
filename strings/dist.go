package strings

// PDistance returns the proportion of bytes that differ in s and t.  If one
// of the strings is longer than the other, the extra bytes are ignored.
func PDistance(s, t string) float64 {
	n := len(s)
	if len(t) < n {
		n = len(t)
	}
	return float64(HammingDistance(s, t)) / float64(n)
}

// Returns the number of bytes that differ in s and t.
// If one of the strings is longer than the other, the extra bytes are ignored.
func HammingDistance(s, t string) int {
	d := 0
	for i := 0; i < len(s) && i < len(t); i++ {
		if s[i] != t[i] {
			d++
		}
	}
	return d
}

func DistanceMatrix(a []string, f func(s, t string) float64) [][]float64 {
	n := len(a)
	b := make([]float64, n * n)
	m := make([][]float64, n)
	for i := range a {
		m[i] = b[n * i:n * (i + 1)]
	}
	for i, s := range a {
		for j := i + 1; j < len(a); j++ {
			t := a[j]
			d := f(s, t)
			m[i][j] = d
			m[j][i] = d
		}
	}
	return m
}
