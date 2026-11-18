package parser

import (
	"databaseInGo/lexer"
	"fmt"
)

func tokenFromKeyword(keyword lexer.Keyword) lexer.Token {
	return lexer.Token{
		Value: string(keyword),
		Type:  lexer.KeywordType,
	}
}

func tokenFromSymbol(symbol lexer.Symbol) lexer.Token {
	return lexer.Token{
		Value: string(symbol),
		Type:  lexer.SymbolType,
	}
}

func isToken(tokens []*lexer.Token, cursor uint, otherToken lexer.Token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return otherToken.Equals(tokens[cursor])
}

func helpMessage(tokens []*lexer.Token, cursor uint, helpMessage string) {
	var c *lexer.Token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}
	fmt.Printf("[%d, %d]: %s, got: %s\n", c.Loc.Line, c.Loc.Col, helpMessage, c.Value)
}

func parseExpression(tokens []*lexer.Token, cursor uint) (*Expression, uint, bool) {
	cur := cursor
	if cur >= uint(len(tokens)) {
		helpMessage(tokens, cur, "out of bounds")
		return nil, cursor, false
	}
	curToken := tokens[cur]
	kinds := []lexer.TokenType{lexer.IdentifierType, lexer.NumType, lexer.TextType}
	for _, kind := range kinds {
		if curToken.Type == kind {
			return &Expression{
					curToken,
					literalExpression,
				},
				cur + 1,
				true
		}
	}
	helpMessage(tokens, cur, "invalid token type")
	return nil, cursor, false
}

func parseExpressions(tokens []*lexer.Token, cursor uint, delimiter []lexer.Token) (*[]*Expression, uint, bool) {
	cur := cursor
	if cur >= uint(len(tokens)) {
		helpMessage(tokens, cur, "out of bounds")
		return nil, cursor, false
	}
	expressions := []*Expression{}
outer:
	for cur < uint(len(tokens)) {
		curToken := tokens[cur]
		for _, delimiterToken := range delimiter {
			if delimiterToken.Equals(curToken) {
				break outer
			}
		}

		if len(expressions) > 0 {
			if !isToken(tokens, cur, tokenFromSymbol(lexer.Comma)) {
				helpMessage(tokens, cur, "expected comma")
				return nil, cursor, false
			}
			cur++
		}
		curExp, newCursor, ok := parseExpression(tokens, cur)
		if !ok {
			return nil, cursor, false
		}
		expressions = append(expressions, curExp)
		cur = newCursor
	}
	return &expressions, cur, true
}

func parseColumnDefs(tokens []*lexer.Token, cursor uint, delimiter []lexer.Token) (*[]*ColumnDefinition, uint, bool) {
	cur := cursor
	if cur >= uint(len(tokens)) {
		return nil, cursor, false
	}

	coldefs := []*ColumnDefinition{}
outer:
	for cur < uint(len(tokens)) {
		curToken := tokens[cur]
		for _, delimiterToken := range delimiter {
			if delimiterToken.Equals(curToken) {
				break outer
			}
		}
		if len(coldefs) > 0 {
			if !isToken(tokens, cur, tokenFromSymbol(lexer.Comma)) {
				helpMessage(tokens, cur, "expected comma")
				return nil, cursor, false
			}
			cur++
		}

		colName, newCursor, ok := parseToken(tokens, cur, lexer.IdentifierType)
		if !ok {
			helpMessage(tokens, cur, "expected identifier")
			return nil, cursor, false
		}
		cur = newCursor
		if cur >= uint(len(tokens)) {
			helpMessage(tokens, cur, "expected more tokens")
			return nil, cursor, false
		}
		colType, newCursor, ok := parseToken(tokens, cur, lexer.KeywordType)
		if !ok {
			helpMessage(tokens, cur, "expected identifier")
			return nil, cursor, false
		}
		cur = newCursor
		coldefs = append(coldefs, &ColumnDefinition{
			ColName: *colName,
			ColType: *colType,
		})
	}
	return &coldefs, cur, true
}

func parseToken(tokens []*lexer.Token, cursor uint, ttype lexer.TokenType) (*lexer.Token, uint, bool) {
	cur := cursor
	if cur >= uint(len(tokens)) {
		return nil, cursor, false
	}
	curToken := tokens[cur]
	if curToken.Type == ttype {
		return curToken, cur + 1, true
	}
	return nil, cursor, false
}
