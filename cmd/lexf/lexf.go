package main

import "bufio"
import "fmt"
import "log"
import "strconv"
import "strings"
import "os"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	fields := strings.Fields(line)

	alphabet := make([]byte, len(fields))
	for i, f := range fields {
		alphabet[i] = f[0]
	}

	line, err = br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	n64, err := strconv.ParseInt(strings.TrimSpace(line), 0, 32)
	if err != nil {
		log.Fatal(err)
	}
	n := int(n64)

	indices := make([]int, n)
	for {
		for _, i := range indices {
			fmt.Printf("%c", alphabet[i])
		}
		fmt.Println()

		i := len(indices) - 1
		for i >= 0 {
			indices[i]++
			if indices[i] == len(alphabet) {
				indices[i] = 0
				i--
			} else {
				break
			}
		}
		if i < 0 {
			break
		}
	}
} 
