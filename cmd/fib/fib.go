package main

import "fmt"

func main() {
	var n int
	var k int64
	fmt.Scan(&n, &k)

	a := int64(1)
	b := int64(1)
	for i := 1; i < n; i++ {
		a, b = b, k * a + b
	}
	fmt.Println(a)
}
