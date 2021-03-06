package tests

import (
	"bytes"
	"testing"

	"github.com/arr-ai/arrai/syntax"
)

func lxr(input string) *syntax.Lexer {
	return syntax.NewLexer(bytes.NewBufferString(input))
}

func TestLexSymbols(t *testing.T) {
	t.Parallel()
	AssertScan(t, lxr(" \t{aa"), syntax.Token('{'), nil, "{")
	AssertScan(t, lxr("\t }|}\n "), syntax.Token('}'), nil, "}")
	AssertScan(t, lxr("\t +12\n "), syntax.Token('+'), nil, "+")
	AssertScan(t, lxr("\t -1\n "), syntax.Token('-'), nil, "-")
	AssertScan(t, lxr("\t ,2\n "), syntax.Token(','), nil, ",")
	AssertScan(t, lxr("\t :{|0\n "), syntax.Token(':'), nil, ":")
	AssertScan(t, lxr("   {||}"), syntax.OSET, nil, "{|")
	AssertScan(t, lxr("   |}}"), syntax.CSET, nil, "|}")
	AssertScan(t, lxr("   "), syntax.EOF, nil, "")
}

func TestLexIdent(t *testing.T) {
	t.Parallel()
	AssertScan(t, lxr(" @,"), syntax.IDENT, nil, "@")
	AssertScan(t, lxr(" a,{12"), syntax.IDENT, nil, "a")
	AssertScan(t, lxr(" Ab|}"), syntax.IDENT, nil, "Ab")
	AssertScan(t, lxr(" \na@b 1"), syntax.IDENT, nil, "a@b")
	AssertScan(t, lxr(" \n\t a@b_123__"), syntax.IDENT, nil, "a@b_123__")
}

func TestLexNumber(t *testing.T) {
	t.Parallel()
	AssertScan(t, lxr(" 0}"), syntax.NUMBER, 0, "0")
	AssertScan(t, lxr(" 123,"), syntax.NUMBER, 123, "123")
	AssertScan(t, lxr(" 0.32 |}"), syntax.NUMBER, 0.32, "0.32")
	AssertScan(t, lxr(" 4.5e+123}"), syntax.NUMBER, 4.5e+123, "4.5e+123")
}

func TestLexString(t *testing.T) {
	t.Parallel()
	AssertScan(t, lxr(" \"\"}"), syntax.STRING, nil, "\"\"")
	AssertScan(t, lxr(" \"abc\","), syntax.STRING, nil, "\"abc\"")
	AssertScan(t, lxr(" \"\\t\\n\"|}"), syntax.STRING, nil, "\"\\t\\n\"")
	AssertScan(t, lxr(" \"\\\"\":"), syntax.STRING, nil, "\"\\\"\"")
}

func TestLexSequence(t *testing.T) {
	t.Parallel()
	l := lxr("1 2 3")
	if AssertScan(t, l, syntax.NUMBER, 1, "1") {
		if AssertScan(t, l, syntax.NUMBER, 2, "2") {
			if AssertScan(t, l, syntax.NUMBER, 3, "3") {
				AssertScan(t, l, syntax.EOF, nil, "")
			}
		}
	}
}

// func TestLexBadInput(t *testing.T) {
// 	t.Parallel()
// 	l := lxr("{|<")
// 	assert.Equal(t, syntax.OSET, l.Scan())
// 	assert.Equal(t, syntax.ERROR, l.Scan())
// 	assert.Equal(t, syntax.ERROR, l.Scan())
// }
