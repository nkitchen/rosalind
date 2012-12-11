package main

import "bufio"
import "fmt"
import "log"
import "regexp"
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

	r := regexp.MustCompile(regexp.QuoteMeta(t))
	loc := r.FindAllStringIndex(s, -1)
	for i := range loc {
		fmt.Print(loc[i][0] + 1, " ")
	}
	fmt.Println()
} 
