package main

import "bufio"
import "fmt"
import "rosalind/gene"
import "io"
import "log"
import "net/http"
import "os"
import "regexp"

func main() {
	br := bufio.NewReader(os.Stdin)

	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		id := line[:len(line) - 1]
		findMotif(id)
	}
} 

var nGlycoMotifRe = regexp.MustCompile( `N[^P][ST][^P]` )

func findMotif(accessId string) {
	resp, err := http.Get("http://www.uniprot.org/uniprot/" + accessId + ".fasta")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	s, err := gene.ReadFasta(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	start := 0
	for {
		loc := nGlycoMotifRe.FindStringIndex(s.Data[start:])
		if loc == nil {
			break
		}

		if start == 0 {
			fmt.Println(accessId)
		}
		fmt.Print(start + loc[0] + 1, " ")
		start += loc[0] + 1
	}
	if start > 0 {
		fmt.Println()
	}
}
