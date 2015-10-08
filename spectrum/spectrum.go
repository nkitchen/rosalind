package spectrum

import "fmt"
import "log"
import "math"
import "sort"
import "strconv"

// Computational mass spectrometry

const debug = false

func dprintf(format string, a ...interface{}) {
	if debug {
		fmt.Printf(format, a...)
	}
}

// Molecular mass in fractions of a dalton (Precision)
type Mass int64

const Precision = 1e-4

const Dalton Mass = 1e4

// Parses a decimal string, e.g., 123.456.
func ParseMass(s string) (Mass, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

    m := math.Floor(f / Precision + 0.5)
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
	return fmt.Sprintf("%.4f", float64(m)/float64(Dalton))
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
	return a == b
}

// Returns a key in MonoisotopicMass whose value is m,
// plus or minus Precision.
// The bool return value indicates whether a match was found.
func ResidueByMass(m Mass) (byte, bool) {
	n := len(sortedByMass)
	i := sort.Search(n, func(i int) bool {
		return m - 1 < sortedByMass[i].mass
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
	// What the spectrum came from (e.g., the protein string)
	Source string
}

func New(masses []Mass) (*Spectrum, error) {
	var spec Spectrum
	spec.masses = make([]Mass, len(masses))
	copy(spec.masses, masses)
	sort.Sort(MassSlice(spec.masses))

	return &spec, nil
}

func FromProtein(pr string) (*Spectrum, error) {
	masses := []Mass{}

	prefix := Mass(0)
	for i := 0; i < len(pr); i++ {
		m, ok := MonoisotopicMass[pr[i]]
		if !ok {
			return nil, fmt.Errorf("Invalid amino acid: %c", pr[i])
		}
		prefix += m
		masses = append(masses, prefix)
	}

    suffix := Mass(0)
	for i := len(pr) - 1; i > 0; i-- {
		m := MonoisotopicMass[pr[i]]
		suffix += m
		masses = append(masses, suffix)
	}

	spec, err := New(masses)
	if err == nil {
		spec.Source = pr
	}
	return spec, err
}

// Returns a protein string matching the mass spectrum
// and a bool indicating success.
func (spec *Spectrum) Protein() (string, error) {
	n := len(spec.masses)
	if n%2 == 1 {
		return "", fmt.Errorf("Odd number of masses")
	}

	if spec.complement == nil {
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
					return "", err
				}
			}
			spec.complement[a] = b
			spec.complement[b] = a
		}
	}

	visited := map[Mass]bool{}
	visited[spec.masses[0]] = true
	visited[spec.masses[len(spec.masses) - 1]] = true
	p, found := spec.findProtein(spec.masses[0], visited, nil)
	if found {
		return string(p), nil
	} else {
		return "", fmt.Errorf("Not found")
	}
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

type Convolution struct {
	elem map[Mass]int
}

var ConvUnit float64 = 1e-4

// Returns the spectral convolution of spA and spB
// (Minkowski difference spA - spB).
// Rounds masses to the nearest Precision daltons.
func (spA *Spectrum) Convolution(spB *Spectrum) *Convolution {
	conv := Convolution{map[Mass]int{}}
	for _, a := range spA.masses {
		for _, b := range spB.masses {
			d := a - b
			da := float64(d) / float64(Dalton) / Precision
			r := Mass(math.Floor(da + 0.5) * Precision) * Dalton
			conv.elem[r] = 1 + conv.elem[r]
		}
	}
	return &conv
}

type multisetElem struct {
	mass Mass
	num int
}

// Returns the maximum multiplicity of conv and the corresponding element.
func (conv *Convolution) Max() (Mass, int) {
	best := -1
	which := Mass(0)
	for m, k := range conv.elem {
		if k > best {
			best = k
			which = m
		}
	}
	if best == -1 {
		log.Fatal("No positive multiplicity")
	}

	return which, best
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
