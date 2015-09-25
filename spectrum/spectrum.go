package spectrum

import "errors"
import "fmt"
import "log"
import "math/big"
import "sort"

// Computational mass spectrometry

// A multiple of daltons
type Mass int64

const Dalton Mass = 1e5

var ratDalton *big.Rat = big.NewRat(1e5, 1)

// Parses a decimal string, e.g., 123.456.
// Rounds to the nearest value of Mass.
func ParseMass(s string) (Mass, error) {
	var da big.Rat
	_, ok := da.SetString(s)
	if !ok {
		return 0, errors.New("Parse error")
	}

    var m big.Rat
	m.Mul(&da, ratDalton)

	// Round.
	m.Add(&m, big.NewRat(1, 2))
	var r big.Int
	r.Div(m.Num(), m.Denom())

	return Mass(r.Int64()), nil
}

func mustParseMass(s string) Mass {
	m, err := ParseMass(s)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func (m Mass) String() string {
	return fmt.Sprintf("%.5f", float64(m) / float64(Dalton))
}

var MonoisotopicMass = map[byte]Mass {
	'A': mustParseMass("71.03711"),
	'C': mustParseMass("103.00919"),
	'D': mustParseMass("115.02694"),
	'E': mustParseMass("129.04259"),
	'F': mustParseMass("147.06841"),
	'G': mustParseMass("57.02146"),
	'H': mustParseMass("137.05891"),
	'I': mustParseMass("113.08406"),
	'K': mustParseMass("128.09496"),
	'L': mustParseMass("113.08406"),
	'M': mustParseMass("131.04049"),
	'N': mustParseMass("114.04293"),
	'P': mustParseMass("97.05276"),
	'Q': mustParseMass("128.05858"),
	'R': mustParseMass("156.10111"),
	'S': mustParseMass("87.03203"),
	'T': mustParseMass("101.04768"),
	'V': mustParseMass("99.06841"),
	'W': mustParseMass("186.07931"),
	'Y': mustParseMass("163.06333"),
}

const Tolerance Mass = 10

func approxEqual(a, b Mass) bool {
	d := a - b
	return -Tolerance <= d && d <= Tolerance
}

// Returns a key in MonoisotopicMass whose value is m,
// plus or minus Tolerance.
// The bool return value indicates whether a match was found.
func ResidueByMass(m Mass) (byte, bool) {
	n := len(sortedByMass)
	i := sort.Search(n, func(i int) bool {
		return m - Tolerance <= sortedByMass[i].mass
	})
	if i < n && approxEqual(m, sortedByMass[i].mass) {
		return sortedByMass[i].residue, true
	} else {
		return 0, false
	}
}

type massEntry struct {
	residue byte
	mass Mass
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
