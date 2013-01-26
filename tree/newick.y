%{

package tree

import "fmt"
import "io"
import "os"
import "strconv"
import "text/scanner"

%}

%union {
	node *Node
	edge Edge
	edges []Edge
	number float64
	text string
}

%type <node>	tree node
%type <text>	label
%type <edge>	edge
%type <edges>	edgeSeq

%token <number>	NUMBER
%token <text>	NAME STRING

%%
tree:
	node ';'
	{
		// Stuffing the result into the lexer is weird,
		// but it's the one reentrant way available.
		yylex.(*newickLexer).Tree = $1
	}

node:
	label
	{
		$$ = &Node{Label: $1, Children: nil}
	}
|
	'(' edgeSeq ')' label
	{
		$$ = &Node{Label: $4, Children: $2}
	}
|
	{
		$$ = &Node{}
	}

label:
	NAME
	{
		$$ = $1
	}
|
	STRING
	{
		$$ = $1
	}

edgeSeq:
	edge {
		$$ = []Edge{$1}
	}
|
	edgeSeq ',' edge
	{
		$$ = append($1, $3)
	}

edge:
	node
	{
		$$ = Edge{$1, 0}
	}
|
	node ':' NUMBER
	{
		$$ = Edge{$1, $3}
	}
%%

type newickLexer struct {
	scanner.Scanner
	Tree *Node
}

type stringError string
func (e stringError) Error() string { return string(e) }

func ReadNewick(r io.Reader) (*Node, error) {
	lexer := newLexer(r)
	yyDebug = 4
	rc := yyParse(lexer)
	return lexer.Tree, stringError(fmt.Sprint(rc))
}

func newLexer(r io.Reader) *newickLexer {
	lexer := newickLexer{}
	lexer.Scanner.Init(r)
	lexer.Scanner.Mode = scanner.ScanIdents | scanner.ScanFloats |
	                     scanner.ScanStrings | scanner.SkipComments
	return &lexer
}

func (lexer *newickLexer) Lex(lval *yySymType) int {
	c := lexer.Peek()
	fmt.Println("Lex: c=", c)
	if c == '\'' {
		token, s := lexer.ScanString(c)
		if token != STRING {
			fmt.Fprintf(os.Stderr, "Error while reading string '%s': %v\n",
			            s, token)
			return 0
		}
		lval.text = s
		return STRING
	}

	token := lexer.Scan()
	fmt.Println("Lex: token=", token)
	switch token {
	case scanner.EOF:
		return 0
	case scanner.Ident:
		lval.text = lexer.TokenText()
		fmt.Println("Lex: Ident", lval.text)
		return NAME
	case scanner.Float:
	    f, err := strconv.ParseFloat(lexer.TokenText(), 64)
		if err != nil {
		    panic(err)
		}
		lval.number = f
		return NUMBER
	case scanner.String:
	    lval.text = lexer.TokenText()
		return STRING
	}

	return int(token)
}

func (lexer *newickLexer) ScanString(quote rune) (int, string) {
	buf := []rune{}
	c := lexer.Next()
	if c != quote {
		panic(c)
	}
	for {
		c = lexer.Next()
		if c == scanner.EOF {
			return int(c), string(buf)
		}
		if c == quote {
			break
		}
		buf = append(buf, c)
	}
	return STRING, string(buf)
}

func (lexer *newickLexer) Error(e string) {
	fmt.Fprintln(os.Stderr, "Error: ", e)
}
