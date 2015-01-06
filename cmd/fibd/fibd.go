package main

import "fmt"

var m int64

func memoize(f func(int64) int64) func(int64) int64 {
	memo := map[int64]int64{}
	return func(x int64) int64 {
		y, ok := memo[x]
		if !ok {
			y = f(x)
			memo[x] = y
		}
		return y
	}
}

var B func(n int64) int64
var M func(n int64) int64

func Br(n int64) int64 {
	switch {
	case n < 0: return 0
	case n == 1: return 1
	default: return M(n - 1)
	}
}

func Mr(n int64) int64 {
	if n < 1 {
		return 0
	}
	return B(n - 1) + M(n - 1) - B(n - m)
}

func F(n int64) int64 {
	return B(n) + M(n)
}

func main() {
	var n int64
	fmt.Scan(&n, &m)

	B = memoize(Br)
	M = memoize(Mr)
	//B = Br
	//M = Mr

    fmt.Println(F(n))
}
