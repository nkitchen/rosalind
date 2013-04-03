package phylo

import "math/big"

type CharArray struct {
	len int
	data *big.Int
}

func NewCharArray(len int) CharArray {
	return CharArray{len, big.NewInt(0)}
}

func (a CharArray) Len() int {
	return a.len
}

func (a CharArray) At(i int) byte {
	return byte(a.data.Bit(i))
}

func (a *CharArray) Set(i int, b byte) {
	a.data.SetBit(a.data, i, uint(b))
}

// PopCount returns the population count of a--the number of 1's in it.
func (a CharArray) PopCount() int {
	n := 0
	for i := 0; i < a.Len(); i++ {
		n += int(a.At(i))
	}
	return n
}

func (c *CharArray) Or(a, b CharArray) {
	if a.Len() != b.Len() {
		panic("Length mismatch")
	}
	if c.Len() != a.Len() {
		*c = NewCharArray(a.Len())
	}

	c.data.Or(a.data, b.data)
}

func (c *CharArray) Not(a CharArray) {
	if c.Len() != a.Len() {
		*c = NewCharArray(a.Len())
	}

	c.data.Not(a.data)
}
