package main

import "bufio"
import "fmt"
import "io"
import "math/big"
import "os"
import "rosalind/tree"
import "strings"

type CharTable []*big.Int

var big1 = big.NewInt(1)

func main() {
	br := bufio.NewReader(os.Stdin)

	line, _ := br.ReadString('\n')
	taxa := strings.Fields(strings.TrimSpace(line))
	// Reverse the order, so that the indices match the bit indices
	// of the character arrays.
	for i, j := 0, len(taxa) - 1; i < j ; {
		taxa[i], taxa[j] = taxa[j], taxa[i]
		i++
		j--
	}

	maskAnchor := big1
	maskAll := &big.Int{}
	maskAll.Lsh(big1, uint(len(taxa)))
	maskAll.Sub(maskAll, big1)

	charTab := CharTable(nil)
	for {
		a := &big.Int{}
		n, err := fmt.Fscanf(br, "%b ", a)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Normalize the array to be 1 at [0].
		if a.Bit(0) == 0 {
			a.Not(a)
			a.And(a, maskAll)
		}
		charTab = append(charTab, a)
	}

	nodeRoot := &tree.Node{}
	nodeAnchor := &tree.Node{Label: taxa[0]}
	nodeRoot.Children = append(nodeRoot.Children, tree.Edge{Node: nodeAnchor})

	maskRest := &big.Int{}
	maskRest.AndNot(maskAll, maskAnchor)

	coanchor := -1
	for i := 1; i < len(taxa); i++ {
		if equalColumns(charTab, 0, i) {
			coanchor = i
			break
		}
	}
	
	if coanchor == -1 {
		i, j := findComplements(charTab, maskRest, maskAnchor)
		if i < 0 || j < 0 {
			panic("No complements found")
		}

		for _, k := range []int{i, j} {
			mask := &big.Int{}
			mask.And(maskRest, charTab[k])
			maskExcl := &big.Int{}
			maskExcl.AndNot(maskAll, mask)
			t := makeSubtree(mask, maskExcl, charTab, taxa)
			nodeRoot.Children = append(nodeRoot.Children, t)
		}
	} else {
		nodeCoanchor := &tree.Node{Label: taxa[coanchor]}
		nodeRoot.Children = append(nodeRoot.Children, tree.Edge{Node: nodeCoanchor})

		maskCoanchor := &big.Int{}
		maskCoanchor.Lsh(big1, uint(coanchor))
		maskRest.AndNot(maskAll, maskCoanchor)
		maskExcl := &big.Int{}
		maskExcl.AndNot(maskAll, maskRest)
		t := makeSubtree(maskRest, maskExcl, charTab, taxa)
		nodeRoot.Children = append(nodeRoot.Children, t)
	}

	tree.Print(nodeRoot)
}

func equalColumns(charTab CharTable, i, j int) bool {
	for k := range charTab {
		if charTab[k].Bit(i) != charTab[k].Bit(j) {
			return false
		}
	}
	return true
}

// findComplements returns the indices of two character arrays a and b
// in charTab such that (a & maskIncl) == (~b & maskIncl),
// a & maskExcl == maskExcl, and b & maskExcl == maskExcl.
func findComplements(charTab CharTable, maskIncl, maskExcl *big.Int) (s, t int) {
	if maskIncl.BitLen() == 0 {
		panic("Empty mask")
	}
	var minInclBit int
	for minInclBit = 0; minInclBit < maskIncl.BitLen(); minInclBit++ {
		if maskIncl.Bit(minInclBit) != 0 {
			break
		}
	}

	arraysByNorm := map[string][]int{}
	for i, a := range charTab {
		b := &big.Int{}
		b.And(a, maskExcl)
		if b.Cmp(maskExcl) != 0 {
			continue
		}

		if a.Bit(minInclBit) == 0 {
			b.AndNot(maskIncl, a)
		} else {
			b.And(maskIncl, a)
		}
		key := fmt.Sprintf("%x", b)
		arraysByNorm[key] = append(arraysByNorm[key], i)
	}
	for _, cands := range arraysByNorm {
		if len(cands) < 2 {
			continue
		}

		for i := 0; i < len(cands) - 1; i++ {
			a := charTab[cands[i]]
			for j := i + 1; j < len(cands); j++ {
				b := charTab[cands[j]]
				c := &big.Int{}
				c.Xor(a, b)
				c.And(c, maskIncl)
				if c.Cmp(maskIncl) == 0 {
					return cands[i], cands[j]
				}
			}
		}
	}
	return -1, -1
}

func findSetBit(s *big.Int, start int) int {
	for i := start; i < s.BitLen(); i++ {
		if s.Bit(i) == 1 {
			return i
		}
	}
	return -1
}

func makeSubtree(maskIncl, maskExcl *big.Int, charTab CharTable,
                 taxa []string) tree.Edge {

	if maskIncl.BitLen() == 0 {
		panic("Empty mask")
	}
	setBit1 := findSetBit(maskIncl, 0)
	setBit2 := findSetBit(maskIncl, setBit1 + 1)
	if setBit2 < 0 {
		e := tree.Edge{Node: &tree.Node{Label: taxa[setBit1]}}
		return e
	}

	setBit3 := findSetBit(maskIncl, setBit2 + 1)
	if setBit3 < 0 {
		f := tree.Edge{Node:&tree.Node{Label: taxa[setBit1]}}
		g := tree.Edge{Node:&tree.Node{Label: taxa[setBit2]}}
		e := tree.Edge{Node:&tree.Node{Children: []tree.Edge{f, g}}}
		return e
	}

	for _, a := range charTab {
		b := &big.Int{}
		b.And(a, maskIncl)
		b.AndNot(b, maskExcl)
		if numSetBits(b) == 1 {
			in := &big.Int{}
			in.AndNot(maskIncl, b)
			out := &big.Int{}
			out.Or(maskExcl, b)
			f := makeSubtree(in, out, charTab, taxa)
			i := b.BitLen() - 1
			g := tree.Edge{Node:&tree.Node{Label: taxa[i]}}
			e := tree.Edge{Node:&tree.Node{Children: []tree.Edge{f, g}}}
			return e
		}
	}

	i, j := findComplements(charTab, maskIncl, maskExcl)
	if i < 0 || j < 0 {
		panic("No complements found")
	}

	n := &tree.Node{}
	for _, k := range []int{i, j} {
		in := &big.Int{}
		in.And(maskIncl, charTab[k])
		notIn := &big.Int{}
		notIn.Not(in)
		out := &big.Int{}
		out.Or(maskExcl, notIn)
		t := makeSubtree(in, out, charTab, taxa)
		n.Children = append(n.Children, t)
	}
	return tree.Edge{Node: n}
}

func numSetBits(n *big.Int) int {
	k := 0
	x := &big.Int{}
	x.Set(n)
	for x.BitLen() > 0 {
		y := &big.Int{}
		y.Sub(x, big1)
		x.And(x, y)
		k++
	}
	return k
}
