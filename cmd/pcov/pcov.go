package main

import "bufio"
import "fmt"
import "os"

func main() {
	next := map[string]string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		n := len(s)
		next[s[:n - 1]] = s[1:]
	}

	super := []byte{}
	visited := map[string]bool{}

	var node string
	for node = range next {
		break
	}
    for {
		node = next[node]
		if visited[node] {
			break
		}

		visited[node] = true
		super = append(super, node[len(node)-1])
	}

	fmt.Println(string(super))
}
