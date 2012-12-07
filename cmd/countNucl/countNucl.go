package main

import "bufio"
import "fmt"
import "os"

func main() {
	counts := make(map[rune]int)
	in := bufio.NewReader(os.Stdin)
	for {
		r, _, err := in.ReadRune()
		if err != nil {
			break
		}
		counts[r] = counts[r] + 1
	}
	for _, k := range ([]rune)("ACGT") {
		fmt.Print(counts[k], " ")
	}
	fmt.Println()
}
