package gene

var DnaComplement = map[byte]byte{'A': 'T', 'C': 'G', 'G': 'C', 'T': 'A'}

func ReverseComplement(dna string) string {
	rc := make([]byte, len(dna))
	for i := 0; i < len(dna); i++ {
		j := len(dna) - 1 - i
		rc[i] = DnaComplement[dna[j]]
	}
	return string(rc)
}
