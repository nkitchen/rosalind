package main

import "bufio"
import "fmt"
import "log"
import "os"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	s := line[:len(line) - 1]

	line, err = br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	t := line[:len(line) - 1]

	fmt.Println(HammingDist(s, t))
}

func HammingDist(s, t string) int {
	if len(s) != len(t) {
		log.Fatalf("Mismatched string lengths: %v vs. %v\n", len(s), len(t))
	}
	d := 0
	for i := range(s) {
		if s[i] != t[i] {
			d++
		}
	}
	return d
}
