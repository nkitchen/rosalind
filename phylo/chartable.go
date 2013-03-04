package phylo

type CharArray []byte

// PopCount returns the population count of a--the number of 1's in it.
func (a CharArray) PopCount() int {
	n := 0
	for _, b := range a {
		n += int(b)
	}
	return n
}

func (c *CharArray) Or(a, b CharArray) {
	if len(a) != len(b) {
		panic("Length mismatch")
	}
	if len(*c) != len(a) {
		*c = make(CharArray, len(a))
	}

	for i := range a {
		(*c)[i] = a[i] | b[i]
	}
}

func (c *CharArray) Not(a CharArray) {
	if len(*c) != len(a) {
		*c = make(CharArray, len(a))
	}

	for i := range a {
		(*c)[i] = 1 - a[i]
	}
}
