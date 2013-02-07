package main

import "bufio"
import "fmt"
import "strings"
import "os"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	_ = err
	taxa := strings.Fields(line)

	seen := map[string]bool{}
	for {
		line, err = br.ReadString('\n')
		if err != nil {
			break
		}

		m := map[byte][]int{}
		s := strings.TrimSpace(line)
		for i := 0; i < len(s); i++ {
			c := s[i]
			m[c] = append(m[c], i)
		}
		c := m['1']
		d := m['0']
		if len(c) < 2 || len(d) < 2 {
			continue
		}

		for i, a1 := range c {
			for j := i + 1; j < len(c); j++ {
				a2 := c[j]

				for k, b1 := range d {
					for l := k + 1; l < len(d); l++ {
						b2 := d[l]

						ta1 := taxa[a1]
						ta2 := taxa[a2]
						if ta2 < ta1 {
							ta1, ta2 = ta2, ta1
						}

						tb1 := taxa[b1]
						tb2 := taxa[b2]
						if tb2 < tb1 {
							tb1, tb2 = tb2, tb1
						}

						if tb1 < ta1 {
							ta1, tb1 = tb1, ta1
							ta2, tb2 = tb2, ta2
						}
						q := fmt.Sprintf("%v,%v:%v,%v", ta1, ta2, tb1, tb2)
						if !seen[q] {
							fmt.Printf("{%v, %v} {%v, %v}\n",
									   ta1, ta2, tb1, tb2)
							seen[q] = true
						}
					}
				}
			}
		}
	}
}
