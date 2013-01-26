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
		fmt.Println("[[", input, "]]")
		t, err := tree.ReadNewick(strings.NewReader(input))
		fmt.Println(err)
		tree.Print(t)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = br.ReadString('\n')

		line, err := br.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		labels := strings.Fields(line)
		fmt.Print(t.Distance(labels[0], labels[1]), " ")

		buf, _ := br.Peek(1)
		if len(buf) > 0 && buf[0] == '\n' {
			_, _ = br.ReadByte()
		}
	}
	fmt.Println()
}
