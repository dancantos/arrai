package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
)

func assertParse(t *testing.T, expected rel.Value, input string) bool { //nolint:unparam
	value, err := syntax.Parse(syntax.NewStringLexer(input))
	return assert.NoError(t, err) &&
		assert.True(t, expected.Equal(value), "%s == \n%s", expected, value)
}

func assertParseError(t *testing.T, input string) bool {
	value, err := syntax.Parse(syntax.NewStringLexer(input))
	return !assert.Error(t, err) &&
		assert.Fail(t, "expected error, got value", "%s", value)
}

func TestParseNumber(t *testing.T) {
	t.Parallel()
	assertParse(t, rel.NewNumber(0), "0")
	assertParse(t, rel.NewNumber(123), "123")
	assertParse(t, rel.NewNumber(0.32), "0.32")
	assertParse(t, rel.NewNumber(4.5e+123), "4.5e+123")
}

func TestParseTuple(t *testing.T) {
	t.Parallel()
	assertParse(t, rel.EmptyTuple, `{}`)
	assertParse(t,
		rel.NewTuple(rel.Attr{Name: "a", Value: rel.NewNumber(1)}),
		`{"a":1}`)
	assertParse(t, rel.NewTuple(
		rel.Attr{Name: "a", Value: rel.NewNumber(1)},
		rel.Attr{Name: "b", Value: rel.NewNumber(2)},
	), `{"a":1, "b": 2}`)
	assertParse(t, rel.NewTuple(
		rel.Attr{Name: "a", Value: rel.NewNumber(1)},
		rel.Attr{Name: "b", Value: rel.NewNumber(2)},
	), `{a :1, b : 2}`)
}

func TestParseSet(t *testing.T) {
	t.Parallel()
	assertParse(t, rel.NewSet(), `{||}`)
	assertParse(t, rel.NewSet(), `false`)
	assertParse(t, rel.NewSet(rel.NewNumber(1)), `{|1|}`)
	assertParse(t, rel.NewSet(rel.NewNumber(1), rel.NewNumber(2)), `{|1,2|}`)
	assertParse(t, rel.NewSet(
		rel.NewNumber(1),
		rel.NewSet(rel.NewNumber(3), rel.NewNumber(4)),
		rel.NewNumber(2),
	), `{|1, {|3, 4|}, 2|}`)
}

func TestParseMixed(t *testing.T) {
	t.Parallel()
	assertParse(t, rel.NewTuple(
		rel.Attr{Name: "a", Value: rel.NewNumber(1)},
		rel.Attr{Name: "b", Value: rel.NewSet(
			rel.NewTuple(rel.Attr{Name: "d", Value: rel.NewNumber(3)}),
			rel.NewNumber(4),
		)},
		rel.Attr{Name: "c", Value: rel.NewNumber(2)},
	), `{a:1, b:{|{d:3}, 4,|}, c:2,}`)
}

func TestParseRelationShortcut(t *testing.T) {
	t.Parallel()
	value, err := syntax.Parse(syntax.NewStringLexer(`{|<a,b> {1, 2}, {3, 4}|}`))
	assert.Error(t, err, "%s", value)
}
