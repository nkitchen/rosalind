package main

import "bufio"
import "fmt"
import "rosalind/gene"
import "io"
import "log"
import "os"
import "regexp"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	dna := strings.TrimSpace(line)

	introns := []string{}
	for {
		line, err = br.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		introns = append(introns, strings.TrimSpace(line))
	}

	spliced := make([]byte, 0, len(dna))
	r := regexp.MustCompile(strings.Join(introns, "|"))
	begin := 0
	for {
		loc := r.FindStringIndex(dna[begin:])
		if loc == nil {
			spliced = append(spliced, dna[begin:]...)
			break
		}
		spliced = append(spliced, dna[begin:begin+loc[0]]...)
		begin += loc[1]
	}

	for i, b := range spliced {
		if b == 'T' {
			spliced[i] = 'U'
		}
	}

	prot, _ := gene.Translate(string(spliced))
	fmt.Println(prot)
}
