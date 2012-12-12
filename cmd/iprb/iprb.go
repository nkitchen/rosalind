package main

import "bufio"
import "flag"
import "fmt"
import "log"
import "strconv"
import "strings"
import "os"
import "math/rand"

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		seed, err := strconv.ParseInt(flag.Arg(0), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		rand.Seed(seed)
	}

	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	fields := strings.Fields(line)

    k, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		log.Fatal(err)
	}

    m, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		log.Fatal(err)
	}

    n, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		log.Fatal(err)
	}

	all := k + m + n

    pDom := 0.0
	pDom += k / all
	pDom += m / all * 0.5 // dominant allele
	pDom += m / all * 0.5 * (k + 0.5 * (m - 1)) / (all - 1) // recessive allele
	pDom += n / all * (k + 0.5 * m) / (all - 1)
	
	fmt.Println(pDom)
} 
