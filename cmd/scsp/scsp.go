package main

import "fmt"
import "log"

func main() {
	var s, t string
	_, err := fmt.Scan(&s, &t)
	if err != nil {
		log.Fatal(err)
	}

    fmt.Println(LongestCommonSupersequence(s, t))
}

func LongestCommonSupersequence(s, t string) string {
	type entry struct {
		len int
		b byte
		// Indices of the predecessor entry
		iPred, jPred int
	}

    // a[m][n].len is the length of the longest common
	// supersequence of s[:m] and t[:n].
	a := make([][]entry, len(s) + 1)
	for i := range a {
		a[i] = make([]entry, len(t) + 1)
	}

    a[0][0].len = 0

	for m := 1; m <= len(s); m++ {
		a[m][0].len = m
		a[m][0].b = s[m - 1]
		a[m][0].iPred = m - 1
	}
	for n := 1; n <= len(t); n++ {
		a[0][n].len = n
		a[0][n].b = t[n - 1]
		a[0][n].jPred = n - 1
	}

	for i := 0; i < len(s); i++ {
		for j := 0; j < len(t); j++ {
			var e entry
			e.len = a[i][j+1].len + 1
			e.b = s[i]
			e.iPred = i
			e.jPred = j + 1

			if a[i+1][j].len + 1 < e.len {
				e.len = a[i+1][j].len + 1
				e.b = t[j]
				e.iPred = i + 1
				e.jPred = j
			}

			if s[i] == t[j] && a[i][j].len + 1 <= e.len {
				e.len = a[i][j].len + 1
				e.b = s[i]
				e.iPred = i
				e.jPred = j
			}

			a[i+1][j+1] = e
		}
	}

    i := len(s)
	j := len(t)
    r := make([]byte, 0, a[i][j].len)
	for a[i][j].len > 0 {
		e := a[i][j]
		r = append(r, e.b)
		i = e.iPred
		j = e.jPred
	}

	i = 0
	j = len(r) - 1
	for i < j {
		r[i], r[j] = r[j], r[i]
		i++
		j--
	}
	return string(r)
}
