package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidInt(t *testing.T) {
	tok, cur, ok := numLexer("123", cursor{})
	assert.True(t, ok)
	assert.Equal(t, &Token{Value: "123", Type: NumType, Loc: Location{}}, tok)
	assert.Equal(t, cursor{position: 3, loc: Location{Line: 0, Col: 3}}, cur)
}
func TestValidDecimal(t *testing.T) {
	tok, cur, ok := numLexer("3.5", cursor{})
	assert.True(t, ok)
	assert.Equal(t, &Token{Value: "3.5", Type: NumType, Loc: Location{}}, tok)
	assert.Equal(t, cursor{position: 3, loc: Location{Line: 0, Col: 3}}, cur)
}

func TestValidTrailingPeriod(t *testing.T) {
	tok, cur, ok := numLexer("4.", cursor{})
	assert.True(t, ok)
	assert.Equal(t, &Token{Value: "4.", Type: NumType, Loc: Location{}}, tok)
	assert.Equal(t, cursor{position: 2, loc: Location{Line: 0, Col: 2}}, cur)
}

func TestValidLeadingPeriod(t *testing.T) {
	tok, cur, ok := numLexer(".001", cursor{})
	assert.True(t, ok)
	assert.Equal(t, &Token{Value: ".001", Type: NumType, Loc: Location{}}, tok)
	assert.Equal(t, cursor{position: 4, loc: Location{Line: 0, Col: 4}}, cur)
}

func TestValidExponent(t *testing.T) {
	tok, cur, ok := numLexer("5e2", cursor{})
	assert.True(t, ok)
	assert.Equal(t, &Token{Value: "5e2", Type: NumType, Loc: Location{}}, tok)
	assert.Equal(t, cursor{position: 3, loc: Location{Line: 0, Col: 3}}, cur)
}

func TestValidDecimalNegativeExponent(t *testing.T) {
	tok, cur, ok := numLexer("1.925e-3", cursor{})
	assert.True(t, ok)
	assert.Equal(t, &Token{Value: "1.925e-3", Type: NumType, Loc: Location{}}, tok)
	assert.Equal(t, cursor{position: 8, loc: Location{Line: 0, Col: 8}}, cur)
}

func TestInvalidNumerics(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"lone period", "."},
		{"exponent with no digits", "1e"},
		{"exponent followed by a letter", "1ex"},
		{"exponent with plus but no digits", "1e+"},
		{"exponent with minus but no digits", "1e-"},
		{"two decimal points", "1.2.3"},
		{"two exponent markers", "1e5e5"},
		{"decimal point after exponent", "1e5.5"},
		{"not a number", "abc"},
		{"leading plus is an operator", "+5"},
		{"leading minus is an operator", "-5"},
		{"empty input", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tok, cur, ok := numLexer(c.input, cursor{})
			assert.False(t, ok)
			assert.Nil(t, tok)
			assert.Equal(t, cursor{}, cur)
		})
	}
}

func TestValidText(t *testing.T) {
	tok, cur, ok := textLexer("'diddyparty'", cursor{})
	assert.True(t, ok)
	assert.Equal(t, &Token{Value: "diddyparty", Type: TextType, Loc: Location{}}, tok)
	assert.Equal(t, cursor{position: 12, loc: Location{Line: 0, Col: 12}}, cur)
}

func TestMissingEndingQuoteText(t *testing.T) {
	tok, cur, ok := textLexer("'diddyparty", cursor{})
	assert.False(t, ok)
	assert.Nil(t, tok, nil)
	assert.Equal(t, cur, cursor{})
}

func TestMissingStartingQuoteText(t *testing.T) {
	tok, cur, ok := textLexer("diddyparty'", cursor{})
	assert.False(t, ok)
	assert.Nil(t, tok, nil)
	assert.Equal(t, cur, cursor{})
}

func TestLexValid(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []*Token
	}{
		{
			name:  "star and from",
			input: "select * from users;",
			want: []*Token{
				{Value: "select", Type: KeywordType, Loc: Location{Line: 0, Col: 0}},
				{Value: "*", Type: SymbolType, Loc: Location{Line: 0, Col: 7}},
				{Value: "from", Type: KeywordType, Loc: Location{Line: 0, Col: 9}},
				{Value: "users", Type: IdentifierType, Loc: Location{Line: 0, Col: 14}},
				{Value: ";", Type: SymbolType, Loc: Location{Line: 0, Col: 19}},
			},
		},
		{
			name:  "keywords are prefixes of identifiers",
			input: "select selected from tables;",
			want: []*Token{
				{Value: "select", Type: KeywordType, Loc: Location{Line: 0, Col: 0}},
				{Value: "selected", Type: IdentifierType, Loc: Location{Line: 0, Col: 7}},
				{Value: "from", Type: KeywordType, Loc: Location{Line: 0, Col: 16}},
				{Value: "tables", Type: IdentifierType, Loc: Location{Line: 0, Col: 21}},
				{Value: ";", Type: SymbolType, Loc: Location{Line: 0, Col: 27}},
			},
		},
		{
			name:  "uppercase keyword",
			input: "SELECT 1;",
			want: []*Token{
				{Value: "select", Type: KeywordType, Loc: Location{Line: 0, Col: 0}},
				{Value: "1", Type: NumType, Loc: Location{Line: 0, Col: 7}},
				{Value: ";", Type: SymbolType, Loc: Location{Line: 0, Col: 8}},
			},
		},
		{
			name:  "numeric forms",
			input: "select 4., .5, 1e-3;",
			want: []*Token{
				{Value: "select", Type: KeywordType, Loc: Location{Line: 0, Col: 0}},
				{Value: "4.", Type: NumType, Loc: Location{Line: 0, Col: 7}},
				{Value: ",", Type: SymbolType, Loc: Location{Line: 0, Col: 9}},
				{Value: ".5", Type: NumType, Loc: Location{Line: 0, Col: 11}},
				{Value: ",", Type: SymbolType, Loc: Location{Line: 0, Col: 13}},
				{Value: "1e-3", Type: NumType, Loc: Location{Line: 0, Col: 15}},
				{Value: ";", Type: SymbolType, Loc: Location{Line: 0, Col: 19}},
			},
		},
		{
			name:  "multiple lines",
			input: "select id\nfrom users;",
			want: []*Token{
				{Value: "select", Type: KeywordType, Loc: Location{Line: 0, Col: 0}},
				{Value: "id", Type: IdentifierType, Loc: Location{Line: 0, Col: 7}},
				{Value: "from", Type: KeywordType, Loc: Location{Line: 1, Col: 0}},
				{Value: "users", Type: IdentifierType, Loc: Location{Line: 1, Col: 5}},
				{Value: ";", Type: SymbolType, Loc: Location{Line: 1, Col: 10}},
			},
		},
		{
			name:  "no whitespace between tokens",
			input: "select(1);",
			want: []*Token{
				{Value: "select", Type: KeywordType, Loc: Location{Line: 0, Col: 0}},
				{Value: "(", Type: SymbolType, Loc: Location{Line: 0, Col: 6}},
				{Value: "1", Type: NumType, Loc: Location{Line: 0, Col: 7}},
				{Value: ")", Type: SymbolType, Loc: Location{Line: 0, Col: 8}},
				{Value: ";", Type: SymbolType, Loc: Location{Line: 0, Col: 9}},
			},
		},
		{
			name:  "leading and trailing whitespace",
			input: "  select 1;  ",
			want: []*Token{
				{Value: "select", Type: KeywordType, Loc: Location{Line: 0, Col: 2}},
				{Value: "1", Type: NumType, Loc: Location{Line: 0, Col: 9}},
				{Value: ";", Type: SymbolType, Loc: Location{Line: 0, Col: 10}},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tokens, err := Lex(c.input)
			assert.NoError(t, err)
			assert.Len(t, tokens, len(c.want))
			assert.Equal(t, c.want, tokens)
		})
	}
}
