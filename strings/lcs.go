// Modified from suffixarray/qsufsort.go in the Go standard library.
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import "bytes"
import "fmt"
import "log"
import "sort"

// Returns the longest common substring of the strings in coll.
// It assumes that none of the strings contain the null byte '\0', 
// or else that any null bytes are escaped somehow (e.g., replaced by "\0\0").
func LongestCommonSubstring(coll []string) string {
	bufSize := 0
	for _, s := range coll {
		bufSize += len(s)
	}
	bufSize += 8 * len(coll)

	// Concatenate the strings into a single buffer,
	// appending a suffix "\0<i>\0" to each one
	// so that no string's suffix can be a suffix of another.
    buf := bytes.NewBuffer(make([]byte, 0, bufSize))

	// Begin and end indices for each string within the buffer
	stringRanges := make([][2]int, len(coll))

	for i, s := range coll {
		stringRanges[i][0] = buf.Len()
		buf.WriteString(s)
		stringRanges[i][1] = buf.Len()

		buf.WriteByte(0)
		_, err := fmt.Fprint(buf, i)
		if err != nil {
			log.Print(err)
			return ""
		}
		buf.WriteByte(0)
	}

	sa := suffixArray(buf.Bytes())

    // For the most recent suffixes from distinct strings that all share
	// a prefix of length k, sharing[k] is the set of string indices.
    var sharing []map[int]bool
	var bestSuf []byte
    var prevSuf []byte
	for _, s := range sa {
		whichString := findContainingRange(s, stringRanges)
		if whichString == -1 {
			continue
		}

		end := stringRanges[whichString][1]
		suf := buf.Bytes()[s:end]
		if prevSuf == nil || suf[0] != prevSuf[0] {
			sharing = make([]map[int]bool, len(suf) + 1)
			for k := 1; k <= len(suf); k++ {
				sharing[k] = map[int]bool{whichString: true}
			}
		} else {
			sharing[1][whichString] = true
			if len(sharing[1]) == len(coll) && len(bestSuf) < 1 {
				bestSuf = suf[:1]
			}
			j := 1
			for ; j < len(suf) && j < len(prevSuf); j++ {
				k := j + 1
				if suf[j] == prevSuf[j] {
					sharing[k][whichString] = true
					if len(sharing[k]) == len(coll) && len(bestSuf) < k {
						bestSuf = suf[:k]
					}
				} else {
					break
				}
			}
			k := j + 1
			sharing = sharing[:k]
			for ; k <= len(suf); k++ {
				sharing = append(sharing, map[int]bool{whichString: true})
			}
		}
		prevSuf = suf
	}
	return string(bestSuf)
}

func findContainingRange(i int, ranges [][2]int) int {
	if len(ranges) == 0 {
		panic("No ranges")
	}

	min := 0
	max := len(ranges) - 1
	for min <= max {
		if i < ranges[min][0] {
			return -1
		}
		if i >= ranges[max][1] {
			return -1
		}

		mid := (min + max) / 2
		if i < ranges[mid][0] {
			max = mid - 1
		} else if i >= ranges[mid][1] {
			min = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

func suffixArray(data []byte) []int {
	// initial sorting by first byte of suffix
	sa := sortedByFirstByte(data)
	if len(sa) < 2 {
		return sa
	}
	// initialize the group lookup table
	// this becomes the inverse of the suffix array when all groups are sorted
	inv := initGroups(sa, data)

	// the index starts 1-ordered
	sufSortable := &suffixSortable{sa: sa, inv: inv, h: 1}

	for sa[0] > -len(sa) { // until all suffixes are one big sorted group
		// The suffixes are h-ordered, make them 2*h-ordered
		pi := 0 // pi is first position of first group
		sl := 0 // sl is negated length of sorted groups
		for pi < len(sa) {
			if s := sa[pi]; s < 0 { // if pi starts sorted group
				pi -= s // skip over sorted group
				sl += s // add negated length to sl
			} else { // if pi starts unsorted group
				if sl != 0 {
					sa[pi+sl] = sl // combine sorted groups before pi
					sl = 0
				}
				pk := inv[s] + 1 // pk-1 is last position of unsorted group
				sufSortable.sa = sa[pi:pk]
				sort.Sort(sufSortable)
				sufSortable.updateGroups(pi)
				pi = pk // next group
			}
		}
		if sl != 0 { // if the array ends with a sorted group
			sa[pi+sl] = sl // combine sorted groups at end of sa
		}

		sufSortable.h *= 2 // double sorted depth
	}

	for i := range sa { // reconstruct suffix array from inverse
		sa[inv[i]] = i
	}
	return sa
}

func sortedByFirstByte(data []byte) []int {
	// total byte counts
	var count [256]int
	for _, b := range data {
		count[b]++
	}
	// make count[b] equal index of first occurence of b in sorted array
	sum := 0
	for b := range count {
		count[b], sum = sum, count[b]+sum
	}
	// iterate through bytes, placing index into the correct spot in sa
	sa := make([]int, len(data))
	for i, b := range data {
		sa[count[b]] = i
		count[b]++
	}
	return sa
}

func initGroups(sa []int, data []byte) []int {
	// label contiguous same-letter groups with the same group number
	inv := make([]int, len(data))
	prevGroup := len(sa) - 1
	groupByte := data[sa[prevGroup]]
	for i := len(sa) - 1; i >= 0; i-- {
		if b := data[sa[i]]; b < groupByte {
			if prevGroup == i+1 {
				sa[i+1] = -1
			}
			groupByte = b
			prevGroup = i
		}
		inv[sa[i]] = prevGroup
		if prevGroup == 0 {
			sa[0] = -1
		}
	}
	// Separate out the final suffix to the start of its group.
	// This is necessary to ensure the suffix "a" is before "aba"
	// when using a potentially unstable sort.
	lastByte := data[len(data)-1]
	s := -1
	for i := range sa {
		if sa[i] >= 0 {
			if data[sa[i]] == lastByte && s == -1 {
				s = i
			}
			if sa[i] == len(sa)-1 {
				sa[i], sa[s] = sa[s], sa[i]
				inv[sa[s]] = s
				sa[s] = -1 // mark it as an isolated sorted group
				break
			}
		}
	}
	return inv
}

type suffixSortable struct {
	sa  []int
	inv []int
	h   int
	buf []int // common scratch space
}

func (x *suffixSortable) Len() int           { return len(x.sa) }
func (x *suffixSortable) Less(i, j int) bool { return x.inv[x.sa[i]+x.h] < x.inv[x.sa[j]+x.h] }
func (x *suffixSortable) Swap(i, j int)      { x.sa[i], x.sa[j] = x.sa[j], x.sa[i] }

func (x *suffixSortable) updateGroups(offset int) {
	bounds := x.buf[0:0]
	group := x.inv[x.sa[0]+x.h]
	for i := 1; i < len(x.sa); i++ {
		if g := x.inv[x.sa[i]+x.h]; g > group {
			bounds = append(bounds, i)
			group = g
		}
	}
	bounds = append(bounds, len(x.sa))
	x.buf = bounds

	// update the group numberings after all new groups are determined
	prev := 0
	for _, b := range bounds {
		for i := prev; i < b; i++ {
			x.inv[x.sa[i]] = offset + b - 1
		}
		if b-prev == 1 {
			x.sa[prev] = -1
		}
		prev = b
	}
}

const (
	match = iota
	skipFirst
	skipSecond
	)

type lcsEntry struct {
	// The length of the longest common subsequence found so far
	len int
	// How to go to the previous entry in the chain
	// One of match, skipFirst, skipSecond
	prevOp int
}

func LongestCommonSubsequence(s, t string) string {
	if len(s) == 0 || len(t) == 0 {
		return ""
	}

	a := make([][]lcsEntry, len(s))
	for i := range a {
		a[i] = make([]lcsEntry, len(t))
	}

	
	len0 := 0
	if s[0] == t[0] {
		len0 = 1
	}
	a[0][0].len = len0

	for i := 1; i < len(s); i++ {
		a[i][0] = lcsEntry{len0, skipFirst}
	}
	for j := 1; j < len(t); j++ {
		a[0][j] = lcsEntry{len0, skipSecond}
	}

	for i := 1; i < len(s); i++ {
		for j := 1; j < len(t); j++ {
			m := -1
			var prev int
			e := a[i - 1][j]
			if e.len > m {
				m = e.len
				prev = skipFirst
			}
			e = a[i][j - 1]
			if e.len > m {
				m = e.len
				prev = skipSecond
			}

			e = a[i - 1][j - 1]
			if s[i] == t[j] && e.len + 1 > m {
				m = e.len + 1
				prev = match
			}
			a[i][j] = lcsEntry{m, prev}
		}
	}

	// Extract the LCS in reverse order in order to avoid deep recursion.
	i := len(s) - 1
	j := len(t) - 1
	lcs := make([]byte, 0, a[i][j].len)
	for i >= 0 {
		switch a[i][j].prevOp {
		case match:
			lcs = append(lcs, s[i])
			i--
			j--
		case skipFirst:
			i--
		case skipSecond:
			j--
		}
	}
	
	i = 0
	j = len(lcs) - 1
	for i < j {
		lcs[i], lcs[j] = lcs[j], lcs[i]
		i++
		j--
	}
	return string(lcs)
}
