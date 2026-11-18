package parser

import (
	"databaseInGo/lexer"
	"fmt"
)

type Ast struct {
	Statements []*Statement
}

func (a *Ast) HasWrites() bool {
	for _, s := range a.Statements {
		if s.Kind == CreateStateKind || s.Kind == InsertStateKind {
			return true
		}
	}
	return false
}

type StatementKind uint

const (
	SelectStateKind StatementKind = iota
	CreateStateKind
	InsertStateKind
)

type Statement struct {
	SelectStatement *SelectStatement
	InsertStatement *InsertStatement
	CreateStatement *CreateStatement
	Kind            StatementKind
}

type SelectStatement struct {
	Items *[]*Expression
	From  *lexer.Token
}

type InsertStatement struct {
	Table  *lexer.Token
	Values *[]*Expression
}

type CreateStatement struct {
	Table      *lexer.Token
	ColumnDefs *[]*ColumnDefinition
}

type ColumnDefinition struct {
	ColName lexer.Token
	ColType lexer.Token
}
type ExpressionKind uint

const (
	literalExpression ExpressionKind = iota
)

type Expression struct {
	Literal *lexer.Token
	Kind    ExpressionKind
}

func parseSelectStatement(tokens []*lexer.Token, cursor uint) (*SelectStatement, uint, bool) {
	cur := cursor
	//Why do some helpers return updated cursor, and other's don't?
	if !isToken(tokens, cur, tokenFromKeyword(lexer.SelectKeyword)) {
		return nil, cursor, false
	}
	cur++
	selectStatement := &SelectStatement{}
	expressions, newCursor, ok := parseExpressions(tokens, cur, []lexer.Token{tokenFromKeyword(lexer.FromKeyword), tokenFromSymbol(lexer.Semicolon)})
	if !ok {
		return nil, cursor, false
	}
	cur = newCursor
	selectStatement.Items = expressions

	if !isToken(tokens, cur, tokenFromKeyword(lexer.FromKeyword)) {
		return nil, cursor, false
	}

	cur++
	curToken, newCursor, ok := parseToken(tokens, cur, lexer.IdentifierType)
	if !ok {
		helpMessage(tokens, cur, "expected identifier token")
		return nil, cursor, false
	}
	selectStatement.From = curToken
	cur = newCursor
	return selectStatement, cur, true
}

// INSERT INTO tableName VALUES (1, 'phil')
func parseInsertStatement(tokens []*lexer.Token, cursor uint) (*InsertStatement, uint, bool) {
	cur := cursor
	if !isToken(tokens, cur, tokenFromKeyword(lexer.InsertKeyword)) {
		return nil, cursor, false
	}
	cur++
	if !isToken(tokens, cur, tokenFromKeyword(lexer.IntoKeyword)) {
		helpMessage(tokens, cur, "expected into keyword")
		return nil, cursor, false
	}
	cur++
	insertStatement := &InsertStatement{}
	curToken, newCursor, ok := parseToken(tokens, cur, lexer.IdentifierType)
	if !ok {
		helpMessage(tokens, cur, "expected identifier token")
		return nil, cursor, false
	}

	insertStatement.Table = curToken
	cur = newCursor

	if !isToken(tokens, cur, tokenFromKeyword(lexer.ValueKeyword)) {
		return nil, cursor, false
	}

	cur++
	if !isToken(tokens, cur, tokenFromSymbol(lexer.LeftParen)) {
		return nil, cursor, false
	}
	cur++
	expressions, newCursor, ok := parseExpressions(tokens, cur, []lexer.Token{tokenFromSymbol(lexer.Semicolon), tokenFromSymbol(lexer.RightParen)})
	if !ok {
		return nil, cursor, false
	}
	insertStatement.Values = expressions
	cur = newCursor
	if !isToken(tokens, cur, tokenFromSymbol(lexer.RightParen)) {
		return nil, cursor, false
	}
	cur++
	return insertStatement, cur, true
}

// CREATE TABLE users (id INT, name TEXT);
func parseCreateStatement(tokens []*lexer.Token, cursor uint) (*CreateStatement, uint, bool) {
	cur := cursor
	if !isToken(tokens, cur, tokenFromKeyword(lexer.CreateKeyword)) {
		return nil, cursor, false
	}
	cur++
	if !isToken(tokens, cur, tokenFromKeyword(lexer.TableKeyword)) {
		helpMessage(tokens, cur, "expected table keyword")
		return nil, cursor, false
	}
	cur++
	createStatement := &CreateStatement{}
	curToken, newCursor, ok := parseToken(tokens, cur, lexer.IdentifierType)
	if !ok {
		helpMessage(tokens, cur, "expected identifier token")
		return nil, cursor, false
	}

	createStatement.Table = curToken
	cur = newCursor

	if !isToken(tokens, cur, tokenFromSymbol(lexer.LeftParen)) {
		return nil, cursor, false
	}
	cur++

	columnDefs, newCursor, ok := parseColumnDefs(tokens, cur, []lexer.Token{tokenFromSymbol(lexer.RightParen), tokenFromSymbol(lexer.Semicolon)})

	if !ok {
		return nil, cursor, false
	}

	createStatement.ColumnDefs = columnDefs
	cur = newCursor
	if !isToken(tokens, cur, tokenFromSymbol(lexer.RightParen)) {
		return nil, cursor, false
	}
	cur++
	return createStatement, cur, true
}

func parseStatement(tokens []*lexer.Token, cursor uint) (*Statement, uint, bool) {
	cur := cursor
	createStatement, newCur, ok := parseCreateStatement(tokens, cur)
	if ok {
		cur = newCur
		return &Statement{
				CreateStatement: createStatement,
				Kind:            CreateStateKind,
			},
			cur,
			true
	}
	selectStatement, newCur, ok := parseSelectStatement(tokens, cur)
	if ok {
		cur = newCur
		return &Statement{
				SelectStatement: selectStatement,
				Kind:            SelectStateKind,
			},
			cur,
			true
	}
	insertStatement, newCur, ok := parseInsertStatement(tokens, cur)
	if ok {
		cur = newCur
		return &Statement{
				InsertStatement: insertStatement,
				Kind:            InsertStateKind,
			},
			cur,
			true
	}
	helpMessage(tokens, cursor, "statement not found")
	return nil, cursor, false
}

func Parse(tokens []*lexer.Token, cursor uint) (*Ast, error) {
	cur := cursor
	ast := Ast{}
	for cur < uint(len(tokens)) {
		statement, newCursor, ok := parseStatement(tokens, cur)
		if !ok {
			return nil, fmt.Errorf("invalid statement starting at token pos %d", cur)
		}
		ast.Statements = append(ast.Statements, statement)
		cur = newCursor

		if cur >= uint(len(tokens)) {
			return nil, fmt.Errorf("Expected a token")
		}
		if !isToken(tokens, cur, tokenFromSymbol(lexer.Semicolon)) {
			return nil, fmt.Errorf("Expected semicolon to end statement, found: %s", tokens[cur].Value)
		}
		cur++
	}
	return &ast, nil
}
