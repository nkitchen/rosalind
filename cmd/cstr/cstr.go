package main

import "bufio"
import "fmt"
import "io"
import "os"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	input := []string{}
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		input = append(input, strings.TrimSpace(line))
	}

	ctab := characterTable(input)
	for _, array := range ctab {
		for _, b := range array {
			fmt.Print(b)
		}
		fmt.Println()
	}
}

func characterTable(coll []string) [][]byte {
	tab := [][]byte{}

	for i := 0; i < len(coll[0]); i++ {
		m := map[byte]int{}
		for _, s := range coll {
			m[s[i]] += 1
		}
		if len(m) < 2 {
			continue
		}
		if len(m) > 2 {
			panic("Not characterizable")
		}

		min := len(coll) + 1
		minByte := byte(0)
		for b, n := range m {
			if n < min {
				min = n
				minByte = b
			}
		}
		if min >= 2 {
			a := make([]byte, len(coll))
			for j, s := range coll {
				if s[i] == minByte {
					a[j] = 1
				} else {
					a[j] = 0
				}
			}

			tab = append(tab, a)
		}
	}
	return tab
}
