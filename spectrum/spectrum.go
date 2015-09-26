package spectrum

import "fmt"
import "log"
import "sort"
import "strconv"

// Computational mass spectrometry

const debug = false

func dprintf(format string, a ...interface{}) {
	if debug {
		fmt.Printf(format, a...)
	}
}

type Mass float64

const Precision = 1e-2

const Dalton Mass = 1

// Parses a decimal string, e.g., 123.456.
func ParseMass(s string) (Mass, error) {
	m, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return Mass(m), nil
}

func mustParseMass(s string) Mass {
	m, err := ParseMass(s)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func (m Mass) String() string {
	return fmt.Sprintf("%.5f", float64(m)/float64(Dalton))
}

var MonoisotopicMass = map[byte]Mass{
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

func ApproxEqual(a, b Mass) bool {
	d := a - b
	return -Precision < d && d < Precision
}

// Returns a key in MonoisotopicMass whose value is m,
// plus or minus Precision.
// The bool return value indicates whether a match was found.
func ResidueByMass(m Mass) (byte, bool) {
	n := len(sortedByMass)
	i := sort.Search(n, func(i int) bool {
		return m-Precision < sortedByMass[i].mass
	})
	if i < n && ApproxEqual(m, sortedByMass[i].mass) {
		return sortedByMass[i].residue, true
	} else {
		return 0, false
	}
}

type massEntry struct {
	residue byte
	mass    Mass
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

type Spectrum struct {
	masses []Mass
	// Maps the mass of each b-ion to that of its y-ion, and vice versa.
	complement map[Mass]Mass
}

func New(masses []Mass) (*Spectrum, error) {
	n := len(masses)

	if n%2 == 1 {
		return nil, fmt.Errorf("Odd number of masses")
	}

	var spec Spectrum
	spec.masses = make([]Mass, n)
	copy(spec.masses, masses)
	sort.Sort(MassSlice(spec.masses))

	spec.complement = map[Mass]Mass{}
	var parent Mass
	for i := 0; 2*i < n; i++ {
		a := spec.masses[i]
		b := spec.masses[n-1-i]
		s := a + b
		if i == 0 {
			parent = s
		} else {
			if !ApproxEqual(parent, s) {
				err := fmt.Errorf("Sum of masses does not match parent: "+
					"%v + %v vs. %v",
					a, b, parent)
				return nil, err
			}
		}
		spec.complement[a] = b
		spec.complement[b] = a
	}

	return &spec, nil
}

// Returns a protein string matching the mass spectrum
// and a bool indicating success.
func (spec *Spectrum) Protein() (string, bool) {
	visited := map[Mass]bool{}
	visited[spec.masses[0]] = true
	visited[spec.masses[len(spec.masses) - 1]] = true
	p, found := spec.findProtein(spec.masses[0], visited, nil)
	return string(p), found
}

func (spec *Spectrum) findProtein(
	from Mass, visited map[Mass]bool, prefix []byte) ([]byte, bool) {
	// visited[m] == true if m has already been used to get a residue
	// in prefix.

    dprintf("%v %v %v\n", from, visited, string(prefix))

	if len(visited) == len(spec.masses) {
		return prefix, true
	}

	ma := from
    for _, mb := range spec.masses {
		if mb < ma || visited[mb] {
			continue
		}
		d := mb - ma
		r, ok := ResidueByMass(d)
		if ok {
			mc := spec.complement[mb]
			visited[mb] = true
			visited[mc] = true
			p, found := spec.findProtein(mb, visited, append(prefix, r))
			if found {
				return p, true
			}
			delete(visited, mb)
			delete(visited, mc)
		}
	}

	return nil, false
}

type MassSlice []Mass

func (s MassSlice) Len() int {
	return len(s)
}

func (s MassSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s MassSlice) Less(i, j int) bool {
	return s[i] < s[j]
}
