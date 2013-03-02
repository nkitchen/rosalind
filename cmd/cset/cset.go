package main

import "bufio"
import "bytes"
import "fmt"
import "io"
import "os"
import "strings"

type CharArray []byte
type CharTable []CharArray

func main() {
	br := bufio.NewReader(os.Stdin)
	
	charTab := CharTable{}
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		b := []byte(strings.TrimSpace(line))
		for i := range b {
			b[i] -= '0'
		}
		charTab = append(charTab, b)
	}

	i := findInconsistentArray(charTab)
	for j := range charTab {
		if j != i {
			fmt.Println(charTab[j])
		}
	}
}

func findInconsistentArray(charTab CharTable) int {
	inconsistencies := make([]int, len(charTab))
	for i := 0; i < len(charTab) - 1; i++ {
		a := charTab[i]
		for j := i + 1; j < len(charTab); j++ {
			b := charTab[j]
			if !consistent(a, b) {
				inconsistencies[i]++
				inconsistencies[j]++
			}
		}
	}
	for i, n := range inconsistencies {
		if n > 1 {
			return i
		}
	}
	for i, n := range inconsistencies {
		if n > 0 {
			return i
		}
	}
	return -1
}

func consistent(a, b CharArray) bool {
	c := &CharArray{}
	aCompl := CharArray{}
	aCompl.Not(a)
	bCompl := CharArray{}
	bCompl.Not(b)

	for _, x := range ([]CharArray{a, aCompl}) {
		for _, y := range ([]CharArray{b, bCompl}) {
			c.And(x, y)
			if c.PopCount() == 0 {
				return true
			}
		}
	}
	return false
}

func (r *CharArray) Not(a CharArray) {
	if len(*r) != len(a) {
		*r = make([]byte, len(a))
	}
	for i := range a {
		(*r)[i] = 1 - a[i]
	}
}

func (r *CharArray) And(a, b CharArray) {
	if len(a) != len(b) {
		panic("Mismatched argument lengths")
	}
	if len(*r) != len(a) {
		*r = make([]byte, len(a))
	}
	for i := range a {
		(*r)[i] = a[i] & b[i]
	}
}

func (a CharArray) PopCount() int {
	n := 0
	for _, b := range a {
		n += int(b)
	}
	return n
}

func (a CharArray) String() string {
	w := &bytes.Buffer{}
	for _, b := range a {
		fmt.Fprintf(w, "%v", b)
	}
	return w.String()
}
