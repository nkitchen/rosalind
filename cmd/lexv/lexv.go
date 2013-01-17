package main

import "bufio"
import "fmt"
import "log"
import "os"
import "sort"
import "strconv"
import "strings"

func main() {
	br := bufio.NewReader(os.Stdin)

	line, err := br.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	alphabet := strings.Fields(line)

	line, _ = br.ReadString('\n')
    n, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		log.Fatal(err)
	}

	perms := []string{}
	max := 1
	for k := 1; k <= n; k++ {
		max *= len(alphabet)
		for i := 0; i < max; i++ {
			s := strconv.FormatInt(int64(i), len(alphabet))
			for len(s) < k {
				s = "0" + s
			}
			perms = append(perms, s)
		}
	}

	sort.StringSlice(perms).Sort()

	for _, s := range perms {
		t := make([]byte, len(s))
		for j := 0; j < len(s); j++ {
			d, err := strconv.ParseInt(string(s[j]), len(alphabet), 32)
			if err != nil {
				log.Fatal(err)
			}
			t[j] = alphabet[int(d)][0]
		}
		fmt.Println(string(t))
	}
}
