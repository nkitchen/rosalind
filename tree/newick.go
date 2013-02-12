
//line newick.y:2

package tree
import __yyfmt__ "fmt"
//line newick.y:3
		
import "fmt"
import "io"
import "strconv"
import "text/scanner"


//line newick.y:12
type yySymType struct {
	yys int
	node *Node
	edge Edge
	edges []Edge
	number float64
	text string
}

const NUMBER = 57346
const NAME = 57347
const STRING = 57348

var yyToknames = []string{
	"NUMBER",
	"NAME",
	"STRING",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 400

//line newick.y:101


type newickReader struct {
	scanner.Scanner
	Tree *Node
	err error
}

type stringError string
func (e stringError) Error() string { return string(e) }

func ReadNewick(r io.Reader) (*Node, error) {
	reader := newReader(r)
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
	case scanner.Int, scanner.Float:
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

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 14
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 21

var yyAct = []int{

	11, 16, 5, 13, 14, 7, 8, 19, 6, 7,
	8, 15, 9, 12, 2, 18, 17, 10, 4, 3,
	1,
}
var yyPact = []int{

	0, -1000, 5, -1000, -1000, -1000, 0, -1000, -1000, -1000,
	-6, -1000, -10, 4, 0, -1000, 3, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 20, 13, 19, 18, 2, 0, 17, 11,
}
var yyR1 = []int{

	0, 1, 2, 2, 3, 4, 5, 5, 5, 7,
	7, 6, 8, 8,
}
var yyR2 = []int{

	0, 2, 1, 1, 1, 4, 1, 1, 0, 1,
	3, 2, 0, 2,
}
var yyChk = []int{

	-1000, -1, -2, -3, -4, -5, 8, 5, 6, 7,
	-7, -6, -2, 9, 10, -8, 11, -5, -6, 4,
}
var yyDef = []int{

	8, -2, 0, 2, 3, 4, 8, 6, 7, 1,
	0, 9, 12, 8, 8, 11, 0, 5, 10, 13,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	8, 9, 3, 3, 10, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 11, 7,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line newick.y:33
		{
			yylex.(*newickReader).Tree = yyS[yypt-1].node
		}
	case 2:
		//line newick.y:39
		{
			yyVAL.node = yyS[yypt-0].node;
		}
	case 3:
		//line newick.y:44
		{
			yyVAL.node = yyS[yypt-0].node;
		}
	case 4:
		//line newick.y:50
		{
			yyVAL.node = &Node{Label: yyS[yypt-0].text, Children: nil}
		}
	case 5:
		//line newick.y:56
		{
			yyVAL.node = &Node{Label: yyS[yypt-0].text, Children: yyS[yypt-2].edges}
		}
	case 6:
		//line newick.y:62
		{
			yyVAL.text = yyS[yypt-0].text
		}
	case 7:
		//line newick.y:67
		{
			yyVAL.text = yyS[yypt-0].text
		}
	case 8:
		//line newick.y:71
		{
			yyVAL.text = ""
		}
	case 9:
		//line newick.y:77
		{
			yyVAL.edges = []Edge{yyS[yypt-0].edge}
		}
	case 10:
		//line newick.y:82
		{
			yyVAL.edges = append(yyS[yypt-2].edges, yyS[yypt-0].edge)
		}
	case 11:
		//line newick.y:88
		{
			yyVAL.edge = Edge{yyS[yypt-1].node, yyS[yypt-0].number}
		}
	case 12:
		//line newick.y:93
		{
			yyVAL.number = 0
		}
	case 13:
		//line newick.y:98
		{
			yyVAL.number = yyS[yypt-0].number
		}
	}
	goto yystack /* stack new state and value */
}
