package main

import "bufio"
import "fmt"
import "io"
import "os"

type DnaString struct {
	Id string
	Bases string
}

func main() {
	prefixes := make(map[string][]DnaString)

	dna := ReadAllFasta(os.Stdin)
	k := 3
	for _, t := range dna {
		prefix := t.Bases[:k]
		prefixes[prefix] = append(prefixes[prefix], t)
	}

	for _, s := range dna {
		suffix := s.Bases[len(s.Bases)-k:]
		for _, t := range prefixes[suffix] {
			if s.Id == t.Id {
				continue
			}
			fmt.Println(s.Id, t.Id)
		}
	}
}

func ReadAllFasta(r io.Reader) []DnaString {
	var dna []DnaString
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if (err != nil && err != io.EOF) {
			return dna
		}

		if (len(line) == 0) {
			return dna
		}

		n := len(line) - 1
		if line[n] == '\n' {
			line = line[:n]
		}

		if (line[0] == '>') {
			dna = append(dna, DnaString{Id: line[1:]})
		} else {
			p := &dna[len(dna)-1].Bases
			*p += line
		}

		if err == io.EOF {
			return dna
		}
	}
	return nil
}
