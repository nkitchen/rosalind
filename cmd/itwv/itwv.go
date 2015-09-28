package main

import "bufio"
import "fmt"
import "log"
import "os"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	words := []string{}
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	s := words[0]
	patterns := words[1:]
	n := len(patterns)

	mat := make([][]int, n)
	for i := range mat {
		mat[i] = make([]int, n)
	}

	for i := range patterns {
		iPatHist := histogram(patterns[i])
	PatternPair:
		for j := range patterns {
			jPatHist := histogram(patterns[j])
			wovenHist := add(iPatHist, jPatHist)

			n := len(patterns[i]) + len(patterns[j])
			windowHist := histogram(s[:n])
			if windowHist == wovenHist &&
				interweavable(patterns[i], patterns[j], s[:n]) {
				mat[i][j] = 1
				continue
			}

			for k := n; k < len(s); k++ {
				out := s[k-n]
				in := s[k]
				(*windowHist.CountOf(out))--
				(*windowHist.CountOf(in))++
				if windowHist == wovenHist &&
					interweavable(patterns[i], patterns[j], s[k-n+1:k+1]) {
					mat[i][j] = 1
					continue PatternPair
				}
			}
		}
	}

    for i := range mat {
		for j := range mat[i] {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Print(mat[i][j])
		}
		fmt.Println()
	}
}

// s and t can be interwoven to make r.
func interweavable(s, t, r string) bool {
	if r[0] != s[0] && r[0] != t[0] {
		return false
	}

	m := len(s)
	n := len(t)
	// w[i][j] is true if s[:i] and t[:j] can be interwoven
	// to make r[:i+j].
	w := make([][]bool, m + 1)
	for i := 0; i <= m; i++ {
		w[i] = make([]bool, n + 1)
	}

	w[0][0] = true
	for i := 0; i < m; i++ {
		if r[i] == s[i] {
			w[i+1][0] = true
		} else {
			break
		}
	}
	for j := 0; j < n; j++ {
		if r[j] == t[j] {
			w[0][j+1] = true
		} else {
			break
		}
	}

    for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if w[i - 1][j] && s[i-1] == r[i + j - 1] {
				w[i][j] = true
			} else if w[i][j - 1] && t[j-1] == r[i + j - 1] {
				w[i][j] = true
			}
		}
	}
	return w[m][n]
}

type hist struct {
	A, C, G, T int
}

func (h *hist) CountOf(b byte) *int {
	switch b {
	case 'A':
		return &(h.A)
	case 'C':
		return &(h.C)
	case 'G':
		return &(h.G)
	case 'T':
		return &(h.T)
	}
	return nil
}

func histogram(s string) hist {
	h := hist{}
	for i := 0; i < len(s); i++ {
		(*h.CountOf(s[i]))++
	}
	return h
}

func add(h1, h2 hist) hist {
	hr := hist{
		A: h1.A + h2.A,
		C: h1.C + h2.C,
		G: h1.G + h2.G,
		T: h1.T + h2.T,
	}
	return hr
}
