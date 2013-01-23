package main

import "bufio"
import "fmt"
import "io"
import "log"
import "os"
import "rosalind/tree"
import "strings"

var _ = fmt.Println

func main() {
	br := bufio.NewReader(os.Stdin)

	for {
		input, err := br.ReadString(';')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		t, err := tree.ReadNewick(strings.NewReader(input))
		if err != nil {
			log.Fatal(err)
		}

		tree.Print(t)

		_, _ = br.ReadString('\n')
	}
}
