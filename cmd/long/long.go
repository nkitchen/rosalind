package main

import "bufio"
import "fmt"
import "io"
import "log"
import "os"

type overlap struct {
	which int // index of the overlapping string (its prefix is the overlap)
	loc int // byte index where the overlap starts
}

func main() {
	br := bufio.NewReader(os.Stdin)

    var coll []string
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		s := line[:len(line) - 1]
		coll = append(coll, s)
	}

	// overlap graph
	g := make([][]overlap, len(coll))
	for i, s := range coll {
		for j, t := range coll {
			if i == j {
				continue
			}

			for k := len(s) / 2 + 1; k < len(s) && k < len(t); k++ {
				if s[len(s)-k:] == t[:k] {
					g[i] = append(g[i], overlap{j, k})
				}
			}
		}
	}

    var p []int
	start := make(map[int]int)
	var postorder func(i, s int)
	postorder = func (i, s int) {
		if t, ok := start[i]; ok {
			if s > t {
				start[i] = s
			}
			return
		}
		start[i] = s
		for _, ov := range g[i] {
			postorder(ov.which, ov.loc)
		}
		p = append(p, i)
	}
	for i := range g {
		postorder(i, 0)
	}

	for i := len(p) - 1; i >= 0; i-- {
		j := p[i]
		fmt.Print(coll[j][start[j]:])
	}
	fmt.Println()
} 
