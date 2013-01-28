package main

import "bufio"
import "fmt"
import "os"
import "strconv"
import "strings"

const M = 1000000

func main() {
	br := bufio.NewReader(os.Stdin)
	line, _ := br.ReadString('\n')
	n, _ := strconv.ParseInt(strings.TrimSpace(line), 10, 64)

    p := int64(1)
	f := int64(2)
	for n > 0 {
		if n & 1 != 0 {
			p = (p * f) % M
		}
		f = (f * f) % M
		n >>= 1
	}
	fmt.Println(p)
}
