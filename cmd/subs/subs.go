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
	from := 0
	for {
		loc := r.FindStringIndex(s[from:])
		if loc == nil {
			break
		}

		fmt.Print(from + loc[0] + 1, " ")
		from += loc[0] + 1
	}
	fmt.Println()
} 
