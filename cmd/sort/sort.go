package main

import "bufio"
import "fmt"
import "io"
import "log"
import "os"
import "rosalind/seq"
import "strings"
import "strconv"

var doc = `
1 2 3 4 5 6 7 8 9 10
1 8 9 3 2 7 6 5 4 10

  [      ]
1 5 4 3 2 6 7 8 9 10
      [             ]
1 5 4 10 9 8 7 6 2 3


      [          ]
1 2 3 9 8 7 6 5 4 10
  [        ]
1 7 8 9 3 2 6 5 4 10
`

func main() {
	br := bufio.NewReader(os.Stdin)

Input:
	for {
		p, err := readPerm(br)
		switch err {
			case io.EOF: break Input
			case nil:
			default: log.Fatal(err)
		}
		q, err := readPerm(br)
		if err != nil {
			log.Fatal(err)
		}

		revs, err := seq.Reversals(p, q)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(len(revs))
		for _, r := range revs {
			fmt.Println(r[0] + 1, r[1])
		}
	}
}

func readPerm(br *bufio.Reader) ([]int, error) {
	line, err := br.ReadString('\n')
	for err == nil && strings.TrimSpace(line) == "" {
		line, err = br.ReadString('\n')
	}

    if err != nil {
		return nil, err
	}

	words := strings.Fields(strings.TrimSpace(line))
	p := make([]int, 0, len(words))
	for _, w := range words {
		n, err := strconv.ParseInt(w, 0, 32)
		if err != nil {
			return p, err
		}
		p = append(p, int(n))
	}
	return p, nil
}
