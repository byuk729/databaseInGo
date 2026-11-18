package lexer

import (
	"fmt"
	"strings"
)

type Keyword string

const (
	SelectKeyword Keyword = "select"
	CreateKeyword Keyword = "create"
	TableKeyword  Keyword = "table"
	InsertKeyword Keyword = "insert"
	IntoKeyword   Keyword = "into"
	FromKeyword   Keyword = "from"
	ValueKeyword  Keyword = "values"
	IntKeyword    Keyword = "int"
	TextKeyword   Keyword = "text"
)

type TokenType uint

const (
	KeywordType TokenType = iota
	IdentifierType
	SymbolType
	TextType
	NumType
)

type Symbol string

const (
	Comma      Symbol = ","
	Semicolon  Symbol = ";"
	Asterick   Symbol = "*"
	LeftParen  Symbol = "("
	RightParen Symbol = ")"
	newline    Symbol = "\n"
	emptySpace Symbol = " "
)

type Location struct {
	Line uint
	Col  uint
}

type cursor struct {
	position uint
	loc      Location
}

type Token struct {
	Value string
	Type  TokenType
	Loc   Location
}

type lexer func(dbString string, curCursor cursor) (*Token, cursor, bool)

func numLexer(curString string, c cursor) (*Token, cursor, bool) {
	cur := c
	presentDecimal := false
	presentExponent := false
	presentExpDigit := false
	presentMantissaDigit := false
	for ; cur.position < uint(len(curString)); cur.position++ {
		curChar := curString[cur.position]
		isDigit := '0' <= curChar && curChar <= '9'
		isExponent := curChar == 'e'
		isDecimal := curChar == '.'

		if isDigit {
			if presentExponent {
				presentExpDigit = true
			} else {
				presentMantissaDigit = true
			}
			cur.loc.Col++
			continue
		}
		if isDecimal {
			if presentDecimal {
				return nil, c, false
			}
			presentDecimal = true
			cur.loc.Col++
			continue
		}
		if isExponent {
			// exponent needs a mantissa before it, e.g. reject "e5" and ".e5"
			if presentExponent || !presentMantissaDigit {
				return nil, c, false
			}
			// no decimals allowed after the exponent
			presentDecimal = true
			presentExponent = true
			nextPos := cur.position + 1
			if nextPos < uint(len(curString)) && (curString[nextPos] == '-' || curString[nextPos] == '+') {
				cur.position++
				cur.loc.Col++
			}
			cur.loc.Col++
			continue
		}
		break
	}

	if cur.position == c.position {
		return nil, c, false
	}
	if !presentMantissaDigit {
		return nil, c, false
	}
	if presentExponent && !presentExpDigit {
		return nil, c, false
	}
	return &Token{
		curString[c.position:cur.position],
		NumType,
		c.loc,
	}, cur, true
}

var keywords = map[string]Keyword{
	"select": SelectKeyword,
	"create": CreateKeyword,
	"table":  TableKeyword,
	"insert": InsertKeyword,
	"into":   IntoKeyword,
	"from":   FromKeyword,
	"values": ValueKeyword,
	"int":    IntKeyword,
	"text":   TextKeyword,
}

func keywordLexer(curString string, c cursor) (*Token, cursor, bool) {
	cur := c
	for ; cur.position < uint(len(curString)); cur.position++ {
		curChar := curString[cur.position]
		isLetter := ('a' <= curChar && curChar <= 'z') || ('A' <= curChar && curChar <= 'Z')
		isDigit := '0' <= curChar && curChar <= '9'
		isUnderScore := curChar == '_'
		if !isLetter && !isDigit && !isUnderScore {
			break
		}
		cur.loc.Col++
	}
	if cur.position == c.position {
		return nil, c, false
	}
	word := strings.ToLower(curString[c.position:cur.position])
	if _, ok := keywords[word]; !ok {
		return nil, c, false
	}
	return &Token{
		word,
		KeywordType,
		c.loc,
	}, cur, true
}

var symbolMap = map[string]Symbol{
	",":  Comma,
	";":  Semicolon,
	"*":  Asterick,
	"(":  LeftParen,
	")":  RightParen,
	"\n": newline,
	" ":  emptySpace,
}

func symbolLexer(curString string, c cursor) (*Token, cursor, bool) {
	cur := c
	curChar := curString[cur.position : cur.position+1]
	sym, ok := symbolMap[curChar]
	if ok {
		cur.position++
		cur.loc.Col++
		switch sym {
		case newline:
			cur.loc.Col = 0
			cur.loc.Line++
			return nil, cur, true
		case emptySpace:
			return nil, cur, true
		default:
			return &Token{
					curChar,
					SymbolType,
					c.loc,
				},
				cur,
				true
		}
	}
	return nil, c, false
}
func textLexer(curString string, c cursor) (*Token, cursor, bool) {
	cur := c
	for ; cur.position < uint(len(curString)); cur.position++ {
		curChar := curString[cur.position]
		if cur.position == c.position {
			if curChar != '\'' {
				return nil, c, false
			} else {
				cur.loc.Col++
				continue
			}
		}
		cur.loc.Col++
		if curChar == '\'' {
			cur.position++
			return &Token{
					curString[c.position+1 : cur.position-1],
					TextType,
					c.loc,
				},
				cur,
				true
		}
	}
	return nil, c, false
}
func identifierLexer(curString string, c cursor) (*Token, cursor, bool) {
	cur := c
	for ; cur.position < uint(len(curString)); cur.position++ {
		curChar := curString[cur.position]
		isLetter := ('a' <= curChar && curChar <= 'z') || ('A' <= curChar && curChar <= 'Z')
		isDigit := '0' <= curChar && curChar <= '9'
		isUnderScore := curChar == '_'
		if !isLetter && !isDigit && !isUnderScore {
			break
		}
		cur.loc.Col++
	}
	if cur.position == c.position {
		return nil, c, false
	}
	word := strings.ToLower(curString[c.position:cur.position])
	return &Token{
		word,
		IdentifierType,
		c.loc,
	}, cur, true
}

func Lex(dbString string) ([]*Token, error) {
	tokens := []*Token{}
	cur := cursor{}
lex:
	for cur.position < uint(len(dbString)) {
		for _, curLexer := range []lexer{numLexer, keywordLexer, symbolLexer, textLexer, identifierLexer} {
			parsedToken, newCursor, ok := curLexer(dbString, cur)
			if ok {
				cur = newCursor
				if parsedToken != nil {
					tokens = append(tokens, parsedToken)
				}
				continue lex
			}
		}
		hint := ""
		if len(tokens) > 0 {
			hint = "after " + tokens[len(tokens)-1].Value
		}
		return nil, fmt.Errorf("error parsing token %s, line %d and column %d", hint, cur.loc.Line, cur.loc.Col)
	}
	return tokens, nil
}
