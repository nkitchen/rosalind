package main

import "bytes"
import "fmt"
import "log"
import "os"

func main() {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	s := buf.Bytes()
	for i, b := range s {
		if b == 'T' {
			s[i] = 'U'
		}
	}
	fmt.Println(string(s))
}
