package main

import "bufio"
import "fmt"
import "log"
import "os"
import "sort"
import "rosalind/seq"
import "strconv"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	nums := []int{}
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, x)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	p := nums[1:]

	lis := seq.LongestIncreasingSubseqInts(p)
	display(lis)

	lds2 := seq.LongestIncreasingSubseqIndex(sort.Reverse(sort.IntSlice(p)))
	{
		r := make([]int, len(lds2))
		for i, ri := range lds2 {
			r[i] = p[ri]
		}
		display(r)
	}
}

func display(a []int) {
	for i, x := range a {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(x)
	}
	fmt.Println()
}
