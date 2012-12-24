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

    values := make([]float64, len(fields))
	for i, f := range fields {
		x, err := strconv.ParseFloat(f, 64)
		if err != nil {
			log.Fatal(err)
		}
		values[i] = x
	}

	nHomDom2 := values[0]
	nHomDomHet := values[1]
	nHomDomHomRec := values[2]
	nHet2 := values[3]
	nHetHomRec := values[4]

    expDom := 0.0

	expDom += 2 * nHomDom2
	expDom += 2 * nHomDomHet
	expDom += 2 * nHomDomHomRec
	expDom += 1.5 * nHet2
	expDom += 1 * nHetHomRec
	
	fmt.Println(expDom)
} 
