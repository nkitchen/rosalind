package main

import "bytes"
import "fmt"
import "log"
import "os"

func main() {
	complBases := map[rune]rune{
		'A': 'T',
		'C': 'G',
		'G': 'C',
		'T': 'A',
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	orig := string(buf.Bytes())
	compl := make([]rune, len(orig))

	for i, b := range orig {
		j := len(compl) - 1 - i
		c, found := complBases[b]
		if found {
			compl[j] = c
		} else {
			compl[j] = b
		}
	}
	fmt.Println(string(compl))
}
