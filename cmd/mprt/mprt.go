package main

import "bufio"
import "fmt"
import "rosalind/gene"
import "io"
import "log"
import "net/http"
import "os"

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

	fmt.Println(s.Data)
}
