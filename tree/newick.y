%{

package tree

import "fmt"
import "io"
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
		yylex.(*newickReader).Tree = $1
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
|
	{
		$$ = ""
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

type newickReader struct {
	scanner.Scanner
	Tree *Node
	err error
}

type stringError string
func (e stringError) Error() string { return string(e) }

func ReadNewick(r io.Reader) (*Node, error) {
	reader := newReader(r)
	yyDebug = 4
	rc := yyParse(reader)
	if rc == 0 {
		return reader.Tree, nil
	}
	return nil, reader.err
}

func newReader(r io.Reader) *newickReader {
	reader := newickReader{}
	reader.Scanner.Init(r)
	reader.Scanner.Mode = scanner.ScanIdents | scanner.ScanFloats |
	                      scanner.ScanStrings | scanner.SkipComments
	reader.Scanner.Error = func(s *scanner.Scanner, msg string) {
		reader.err = stringError(msg)
	}
	return &reader
}

func (r *newickReader) Lex(lval *yySymType) int {
	c := r.Peek()
	if c == '\'' {
		token, s := r.ScanString(c)
		if token != STRING {
			msg := fmt.Sprintf("Error while reading string '%s': %v\n",
			             	   s, token)
			r.err = stringError(msg)
			return 0
		}
		lval.text = s
		return STRING
	}

	token := r.Scan()
	switch token {
	case scanner.EOF:
		return 0
	case scanner.Ident:
		lval.text = r.TokenText()
		return NAME
	case scanner.Float:
	    f, err := strconv.ParseFloat(r.TokenText(), 64)
		if err != nil {
		    panic(err)
		}
		lval.number = f
		return NUMBER
	case scanner.String:
	    lval.text = r.TokenText()
		return STRING
	}

	return int(token)
}

func (r *newickReader) ScanString(quote rune) (int, string) {
	buf := []rune{}
	c := r.Next()
	if c != quote {
		panic(c)
	}
	for {
		c = r.Next()
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

func (r *newickReader) Error(e string) {
	pos := r.Pos()
	r.err = stringError(fmt.Sprintf("%s:%v:%v: %s", pos.Filename,
	                                pos.Line, pos.Column, e))
}
