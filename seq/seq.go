package seq

import "sort"

// Longest increasing subsequence of ints
func LongestIncreasingSubseqInts(s []int) []int {
	lisi := LongestIncreasingSubseqIndex(sort.IntSlice(s))
	return IntsAt(s, lisi)
}

// Finds the longest increasing subsequence of data
// and returns of that subsequence.
func LongestIncreasingSubseqIndex(data sort.Interface) []int {
	if data.Len() == 0 {
		return nil
	}

	// minLastIndex[k] is the index of the minimum element of data
	// such that there is an increasing subsequence of length k + 1.
	minLastIndex := make([]int, 1, data.Len())
	// pred[i] is the index of the predecessor of data[i]
	// in an increasing subsequence.
	pred := make([]int, data.Len())

	minLastIndex[0] = 0
	pred[0] = -1
	for i := 1; i < data.Len(); i++ {
		// Find the position of data[i] in the sort order defined by
		// minLastIndex.
		k := sort.Search(len(minLastIndex),
			func(k int) bool {
				return data.Less(i, minLastIndex[k])
			})
		if k == len(minLastIndex) {
			minLastIndex = append(minLastIndex, i)
		} else {
			minLastIndex[k] = i
		}
		if k == 0 {
			pred[i] = -1
		} else {
			pred[i] = minLastIndex[k-1]
		}
	}

	ri := make([]int, len(minLastIndex))
	k := len(minLastIndex) - 1
	j := minLastIndex[k]
	for k >= 0 {
		ri[k] = j
		k--
		j = pred[j]
	}
	return ri

	//_ := `
	//5 1 4 2 3

	//minLastIndex[1] = 0
	//s[0] == 5
	//pred[0] = -1

	//s[1] == 1
	//1 >= 5? no
	//i < 5? yes
	//minLastIndex[1] = 1

	//s[2] == 4
	//4 >= 1? yes
	//minLastIndex[2] = 2

	//s[3] == 2
	//2 >= s[2] == 4? no
	//2 >= s[1] == 1? yes
	//minLastIndex[2] = 3

	//s[4] == 3
	//3 >= s[3] == 2? yes
	//minLastIndex[3] = 4

	//E A D B C AA F

	//s[0] == E
	//minLastIndex[1] = 0
	//pred[0] = -1

	//s[1] == A
	//_A_ < E
	//minLastIndex[1] = 1
	//pred[1] = -1

	//s[2] == D
	//A < _D_
	//minLastIndex[2] = 2
	//pred[2] = minLastIndex[1] == 1

	//s[3] == B
	//A < _B_ < D
	//minLastIndex[2] = 3
	//pred[3] = minLastIndex[1] == 1

	//s[4] == C
	//A < B < _C_
	//minLastIndex[3] = 4
	//pred[4] = minLastIndex[2] == 3

	//s[5] == AA
	//A < _AA_ < B < C
	//minLastIndex[2] = 5
	//pred[5] = minLastIndex[1] == 1

	//s[6] == F
	//A < AA < C < F
	//minLastIndex[4] = 6
	//pred[6] = minLastIndex[3] == 4

	//r = s[rev(6, 4, 3, 1)] == [A, B, C, F]
	//`
}

// IntsAt returns the elements of s at the indices in ii.
func IntsAt(s []int, ii []int) []int {
	a := make([]int, len(ii))
	for j, i := range ii {
		a[j] = s[i]
	}
	return a
}

func SlowLGIS(s []int) []int {
	hasSucc := make([]bool, len(s))
	preds := make([][]int, len(s))

	// Find successors.
	for i := len(s) - 1; i >= 0; i-- {
		for j := i - 1; j >= 0; j-- {
			if s[j] < s[i] {
				preds[i] = append(preds[i], j)
				hasSucc[j] = true
			}
		}
	}

	best := []int{}
	for i := range s {
		if hasSucc[i] {
			continue
		}
		r := longestPath(i, preds)
		if len(r) > len(best) {
			best = r
		}
	}
	return IntsAt(s, best)
}

func longestPath(i int, preds [][]int) []int {
	best := []int(nil)
	for _, j := range preds[i] {
		p := longestPath(j, preds)
		if len(p) > len(best) {
			best = p
		}
	}
	best = append(best, i)
	return best
}
