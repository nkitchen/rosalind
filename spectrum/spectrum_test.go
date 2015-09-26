package spectrum_test

import "math/rand"
import "rosalind/spectrum"
import "sort"
import "testing"

func reversed(s string) string {
	b := []byte(s)
	i := 0
	j := len(b) - 1
	for i < j {
		b[i], b[j] = b[j], b[i]
		i++
		j--
	}
	return string(b)
}

const MaxLen = 20
const MinLen = MaxLen / 2

func TestFullSpectrum(t *testing.T) {
	rand.Seed(1)

	aa := []byte{}
	for a := range spectrum.MonoisotopicMass {
		if a == 'L' {
			continue
		}
		aa = append(aa, a)
	}

	for k := 0; k < 100; k++ {
		n := MinLen + rand.Intn(MaxLen - MinLen + 1)
		prot := make([]byte, n)
		for i := 0; i < n; i++ {
			j := rand.Intn(len(aa))
			prot[i] = aa[j]
		}
		protStr := string(prot)

		w1 := spectrum.Mass(500 + 500 * rand.Float64())
		w2 := spectrum.Mass(500 + 500 * rand.Float64())

		masses := []spectrum.Mass{}
		for i := 0; i <= len(prot); i++ {
			m1 := w1
			for j := 0; j < i; j++ {
				m1 += spectrum.MonoisotopicMass[prot[j]]
			}
			m2 := w2
			for j := i; j < len(prot); j++ {
				m2 += spectrum.MonoisotopicMass[prot[j]]
			}
			masses = append(masses, m1, m2)
		}

		// Shuffle.
		for i := 0; i < len(masses) - 1; i++ {
			j := rand.Intn(len(masses) - i)
			masses[i], masses[i + j] = masses[i + j], masses[i]
		}

        sortedMasses := make([]spectrum.Mass, len(masses))
		copy(sortedMasses, masses)
		sort.Sort(spectrum.MassSlice(sortedMasses))

		spec, err := spectrum.New(masses)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		q, ok := spec.Protein()
		if !ok {
			t.Errorf("No protein found for %v; expected %s", masses, protStr)
		} else if q != protStr && q != reversed(protStr) {
			t.Errorf("Found protein %s; expected %s / %v", q, protStr, masses)
		}
	}
}
