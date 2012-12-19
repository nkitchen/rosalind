// Modified from suffixarray/qsufsort.go in the Go standard library.
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import "sort"

type suffix {
	// index of the string whose suffix this is
	of int
	// index of the first byte in the suffix
	start int
}

// LCS returns the longest common substring of the strings in s.
func LCS(s []string) string {
	sa := suffixArray(s)
	return s[sa[0].of][sa[0].start:]
}

func suffixArray(s []string) []suffix {
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

func sortedByFirstByte(data []string) []suffix {
	// total byte counts
	var count [256]int
	for _, s := range data {
		for _, b := range s {
			count[b]++
		}
	}
	// make count[b] equal index of first occurence of b in sorted array
	sum := 0
	for b := range count {
		count[b], sum = sum, count[b]+sum
	}
	// iterate through bytes, placing index into the correct spot in sa
	totalLen := 0
	for _, s := range data {
		totalLen += len(s)
	}
	sa := make([]int, totalLen)
	for i, s := range data {
		for j, b := range s {
			sa[count[b]] = suffix{i, j}
			count[b]++
		}
	}
	for i, b := range data {
		sa[count[b]] = i
		count[b]++
	}
	return sa
}

func initGroups(sa []suffix, data []string) []int {
	a := make([]int, len(sa))
	inv := make([][]int, len(data))
	for i, s := range data {
		inv[i] = a[:len(s)]
		a = a[len(s):]
	}

	// Separate out the final suffix of each string to the start of its group.
	// This is necessary to ensure the suffix "a" is before "aba"
	// when using a potentially unstable sort.
	groupByte := '\0'
	// first non-final suffix
	nf1 := -1
	for i, suf := range sa {
		b := data[suf.of][suf.start]
		final := suf.start == len(data[suf.of]) - 1
		if i == 0 || b != groupByte {
			groupByte = b
			if !final {
				nf1 = i
			} else {
				nf1 = -1
			}
		} else if final && nf1 >= 0 {
			sa[i], sa[nf1] = sa[nf1], sa[i]
			inv[sa[nf1]] = nf1
			sa[nf1] = -1 // mark it as an isolated sorted group

			j := nf1 + 1
			nf1 = -1
			for ; j < i; j++ {
				s := sa[j]
				if data[s.of][s.start] != groupByte {
					break
				}
				if s.start < len(data[s.of]) - 1 {
					nf1 = j
					break
				}
			}
		}
	}

    // TODO: Fix to handle suffixes at -1
	// label contiguous same-letter groups with the same group number
	prevGroup := len(sa) - 1
	suf := sa[prevGroup]
	groupByte := data[suf.of][suf.start]
	for i := len(sa) - 1; i >= 0; i-- {
		suf = sa[i]
		if b := data[suf.of][suf.start]; b < groupByte {
			if prevGroup == i+1 {
				sa[i+1] = suffix{-1, -1}
			}
			groupByte = b
			prevGroup = i
		}
		inv[suf.of][suf.start] = prevGroup
		if prevGroup == 0 {
			sa[0] = &suffix{-1, -1}
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
