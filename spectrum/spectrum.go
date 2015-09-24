package spectrum

import "math"
import "sort"

// Computational mass spectrometry

var MonoisotopicMass = map[byte]float64 {
	'A': 71.03711,
	'C': 103.00919,
	'D': 115.02694,
	'E': 129.04259,
	'F': 147.06841,
	'G': 57.02146,
	'H': 137.05891,
	'I': 113.08406,
	'K': 128.09496,
	'L': 113.08406,
	'M': 131.04049,
	'N': 114.04293,
	'P': 97.05276,
	'Q': 128.05858,
	'R': 156.10111,
	'S': 87.03203,
	'T': 101.04768,
	'V': 99.06841,
	'W': 186.07931,
	'Y': 163.06333,
}

const MassTolerance = 1e-4

func approxEqual(a, b float64) bool {
	return math.Abs(a - b) <= MassTolerance
}

// Returns a key in MonoisotopicMass whose value is mass,
// plus or minus MassTolerance.
// The bool return value indicates whether a match was found.
func ResidueByMass(mass float64) (byte, bool) {
	n := len(sortedByMass)
	i := sort.Search(n, func(i int) bool {
		return mass - MassTolerance <= sortedByMass[i].mass
	})
	if i < n && approxEqual(mass, sortedByMass[i].mass) {
		return sortedByMass[i].residue, true
	} else {
		return 0, false
	}
}

type massEntry struct {
	residue byte
	mass float64
}

var sortedByMass []massEntry

func init() {
	sortedByMass = make([]massEntry, 0, len(MonoisotopicMass))
	for aa, m := range MonoisotopicMass {
		sortedByMass = append(sortedByMass, massEntry{aa, m})
	}

	sort.Sort(massEntrySlice(sortedByMass))
}

type massEntrySlice []massEntry

func (s massEntrySlice) Len() int {
	return len(s)
}

func (s massEntrySlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s massEntrySlice) Less(i, j int) bool {
	a := s[i]
	b := s[j]
	switch {
	case a.mass < b.mass:
		return true
	case a.mass == b.mass && a.residue < b.residue:
		return true
	}
	return false
}
