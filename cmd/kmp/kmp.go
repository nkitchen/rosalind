package main

import "bufio"
import "flag"
import "fmt"
import "os"
import "reflect"
import "strings"

func main() {
	verify := flag.Bool("verify", false, "Check result using brute force")
	flag.Parse()

	br := bufio.NewReader(os.Stdin)

	line, _ := br.ReadString('\n')
	dna := strings.TrimSpace(line)

	fail := make([]int, len(dna))
	for i := 1; i < len(dna); i++ {
		k := fail[i - 1]
		switch {
		case dna[i] == dna[k]:
			fail[i] = k + 1
		case k > 0:
			k = fail[k - 1]
			if dna[i] == dna[k] {
				fail[i] = k + 1
				break
			}
			fallthrough
		default:
			if dna[i] == dna[0] {
				fail[i] = 1
			}
		}
	}

	if *verify {
		brute := make([]int, len(dna))
		for i := 1; i < len(dna); i++ {
			for k := i; k >= 1; k-- {
				if dna[i+1-k:i+1] == dna[:k] {
					brute[i] = k
					break
				}
			}
		}
		
		if !reflect.DeepEqual(fail, brute) {
			fmt.Fprintln(os.Stderr, "Mismatched results")
			fmt.Fprintln(os.Stderr, fail)
			fmt.Fprintln(os.Stderr, brute)
			os.Exit(1)
		}
	}

	output := fmt.Sprint(fail)
	fmt.Println(strings.Trim(output, "[]"))
}
