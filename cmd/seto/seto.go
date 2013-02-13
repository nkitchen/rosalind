package main

import "bufio"
import "math/big"
import "strings"
import "strconv"
import "fmt"
import "os"
import "io"
import "text/scanner"

var n int

func main() {
	br := bufio.NewReader(os.Stdin)
	
	line, _ := br.ReadString('\n')
	n, _ = strconv.Atoi(strings.TrimSpace(line))

	line, _ = br.ReadString('\n')
	a := readSet(strings.NewReader(strings.TrimSpace(line)))
	line, _ = br.ReadString('\n')
	b := readSet(strings.NewReader(strings.TrimSpace(line)))
	
	c := big.NewInt(0)

	c.Or(a, b)
	printSet(c)

	c.And(a, b)
	printSet(c)

	c.AndNot(a, b)
	printSet(c)
	
	c.AndNot(b, a)
	printSet(c)

	c.Not(a)
	printSet(c)

	c.Not(b)
	printSet(c)
}

func readSet(r io.Reader) *big.Int {
	var s scanner.Scanner
	ret := big.NewInt(0)
	s.Init(r)
	token := s.Scan()
	for token != scanner.EOF {
		switch token {
		case scanner.Int:
			x, _ := strconv.Atoi(s.TokenText())
			bit := big.NewInt(1)
			bit.Lsh(bit, uint(x))
			ret.Or(ret, bit)
		}
		token = s.Scan()
	}
	return ret
}

func printSet(s *big.Int) {
	fmt.Print("{")
	a := []string{}
	for i := 1; i <= n; i++ {
		if s.Bit(i) != 0 {
			a = append(a, strconv.Itoa(i))
		}
	}
	fmt.Print(strings.Join(a, ", "))
	fmt.Println("}")
}
