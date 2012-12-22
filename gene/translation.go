package gene

import "fmt"

type CodonError string 

func (e CodonError) Error() string {
	return string(e)
}

// Translate returns the protein string corresponding to the RNA sequence.
func Translate(rna string) (protein string, err error) {
	buf := make([]byte, 0, len(rna) / 3)
	for i := 0; i < len(rna); i += 3 {
		c := rna[i:i+3]
		a, ok := rnaTable[c]
		if ok {
			if a == '.' {
				break
			}
			buf = append(buf, a)
		} else {
			err = CodonError(fmt.Sprintf("Bad codon %s at position %d", c, i))
			break
		}
	}
	protein = string(buf)
	return
}

// rnaTable[c] is the code for the amino acid coded by the RNA codon c.
var rnaTable = map[string]byte{
	"UUU": 'F',  "CUU": 'L',  "AUU": 'I',  "GUU": 'V',
	"UUC": 'F',  "CUC": 'L',  "AUC": 'I',  "GUC": 'V',
	"UUA": 'L',  "CUA": 'L',  "AUA": 'I',  "GUA": 'V',
	"UUG": 'L',  "CUG": 'L',  "AUG": 'M',  "GUG": 'V',
	"UCU": 'S',  "CCU": 'P',  "ACU": 'T',  "GCU": 'A',
	"UCC": 'S',  "CCC": 'P',  "ACC": 'T',  "GCC": 'A',
	"UCA": 'S',  "CCA": 'P',  "ACA": 'T',  "GCA": 'A',
	"UCG": 'S',  "CCG": 'P',  "ACG": 'T',  "GCG": 'A',
	"UAU": 'Y',  "CAU": 'H',  "AAU": 'N',  "GAU": 'D',
	"UAC": 'Y',  "CAC": 'H',  "AAC": 'N',  "GAC": 'D',
	"UAA": '.',  "CAA": 'Q',  "AAA": 'K',  "GAA": 'E',
	"UAG": '.',  "CAG": 'Q',  "AAG": 'K',  "GAG": 'E',
	"UGU": 'C',  "CGU": 'R',  "AGU": 'S',  "GGU": 'G',
	"UGC": 'C',  "CGC": 'R',  "AGC": 'S',  "GGC": 'G',
	"UGA": '.',  "CGA": 'R',  "AGA": 'R',  "GGA": 'G',
	"UGG": 'W',  "CGG": 'R',  "AGG": 'R',  "GGG": 'G',
}
