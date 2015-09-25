package main

import "bufio"
import "fmt"
import "log"
import "math/big"
import "os"
import "strings"

func main() {
	spec := [][]*big.Rat{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		s := []*big.Rat{}
		for _, w := range strings.Split(line, " ") {
			r := new(big.Rat)
			if _, ok := r.SetString(w); !ok {
				log.Fatalf("Parse error: %s", w)
			}
			s = append(s, r)
		}
		spec = append(spec, s)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

    conv := map[string]int{}
	for _, x := range spec[0] {
		for _, y := range spec[1] {
			var d big.Rat
			s := d.Sub(x, y).String()
			conv[s] = 1 + conv[s]
		}
	}

	max := 0
	where := ""
	for s, m := range conv {
		if m > max {
			max = m
			where = s
		}
	}

    var shift big.Rat
	fmt.Println(max)
	shift.SetString(where)
	fmt.Println(shift.FloatString(5))
}
