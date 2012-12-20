package main

import "bufio"
import "fmt"
import "io"
import "log"
import "os"
import "rosalind/strings"

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

    lcs := strings.LCS(coll)
	fmt.Println(lcs)
} 
