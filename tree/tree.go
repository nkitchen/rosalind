package tree

import "fmt"
import "io"

type Node struct {
	Label string
	Children []*Node
}

type treeError string
func (e treeError) Error() string {
	return string(e)
}

func ReadNewick(r io.ByteScanner) (*Node, error) {
	t, err := readNode(r)
	if err != nil {
		return t, err
	}
	c, err := r.ReadByte()
	switch {
	case err != nil:
		return t, err
	case c == ';':
		return t, nil
	default:
		return t, treeError(fmt.Sprintf("Expected ';' but read '%c'", c))
	}
	return nil, treeError("Reached an unreachable line of code")
}

func readNode(r io.ByteScanner) (*Node, error) {
	c, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	t := &Node{}
	if c == '(' {
ReadChildren:
		for {
			u, err := readNode(r)
			if err != nil {
				return nil, err
			}
			t.Children = append(t.Children, u)
			c, err = r.ReadByte()
			if err != nil {
				return nil, err
			}
			switch c {
			case ')':
				break ReadChildren
			case ',':
			    c, err = r.ReadByte()
				if err != nil {
					return nil, err
				}
				if c != ' ' {
					err = r.UnreadByte()
					if err != nil {
						return nil, err
					}
				}
			default:
				s := fmt.Sprintf("Expected ',' or ')' but read '%c'", c)
				return nil, treeError(s)
			}
		}
	} else {
		err = r.UnreadByte()
		if err != nil {
			return nil, err
		}
	}

	label := []byte{}
ReadLabel:
	for {
		c, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
		switch c {
		case ')', ',', ';':
		    err = r.UnreadByte()
			if err != nil {
				return nil, err
			}
			t.Label = string(label)
			break ReadLabel
		default:
			label = append(label, c)
		}
	}

	return t, nil
}

func Print(t *Node) {
	printSubtree(t, "")
}

func printSubtree(t *Node, prefix string) {
	_ = `
	[a]-+-[b]
	    +-[c]-+-[d]-+-[d1]
		|     |     +-[d2]
		|     |     \-[d3]
		|     +-[e]-+-[e1]
		|     |     +-[e2]
		|     |     \-[e3]
		|     \-[f]
		\-[g]---[h]
	`
	fmt.Print("[", t.Label, "]")

	buf := make([]byte, len(t.Label))
	for i := range buf {
		buf[i] = ' '
	}
	labelSpace := string(buf)

	switch {
	case len(t.Children) == 0:
	    fmt.Println()
	case len(t.Children) == 1:
	    fmt.Print("---")
	    printSubtree(t.Children[0], prefix + " " + labelSpace + "    ")
    default:
		p := prefix + " " + labelSpace + "  | "
	    fmt.Print("-+-")
	    printSubtree(t.Children[0], p)
		var i int
		for i = 1; i < len(t.Children) - 1; i++ {
			fmt.Print(prefix, " ", labelSpace, "  +-")
			printSubtree(t.Children[i], p)
		}
		fmt.Print(prefix, " ", labelSpace, `  \-`)
		printSubtree(t.Children[i], p)
	}
}
