package parser

import (
	"databaseInGo/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidSelectQuery(t *testing.T) {
	tokens := []*lexer.Token{
		{Value: "select", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 0}},
		{Value: "id", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 0, Col: 7}},
		{Value: "from", Type: lexer.KeywordType, Loc: lexer.Location{Line: 1, Col: 0}},
		{Value: "users", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 1, Col: 5}},
		{Value: ";", Type: lexer.SymbolType, Loc: lexer.Location{Line: 1, Col: 10}},
	}

	want := &SelectStatement{
		Items: &[]*Expression{
			{Literal: tokens[1], Kind: literalExpression},
		},
		From: tokens[3],
	}

	got, newCursor, ok := parseSelectStatement(tokens, 0)

	assert.True(t, ok)
	assert.Equal(t, uint(4), newCursor)
	assert.Equal(t, got, want)
}

func TestInsertQuery(t *testing.T) {
	tokens := []*lexer.Token{
		{Value: "insert", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 0}},
		{Value: "into", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 7}},
		{Value: "tablename", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 1, Col: 0}},
		{Value: "values", Type: lexer.KeywordType, Loc: lexer.Location{Line: 1, Col: 10}},
		{Value: "(", Type: lexer.SymbolType, Loc: lexer.Location{Line: 2, Col: 0}},
		{Value: "1", Type: lexer.NumType, Loc: lexer.Location{Line: 3, Col: 0}},
		{Value: ",", Type: lexer.SymbolType, Loc: lexer.Location{Line: 4, Col: 0}},
		{Value: "2", Type: lexer.NumType, Loc: lexer.Location{Line: 5, Col: 0}},
		{Value: ")", Type: lexer.SymbolType, Loc: lexer.Location{Line: 6, Col: 0}},
		{Value: ";", Type: lexer.SymbolType, Loc: lexer.Location{Line: 7, Col: 0}},
	}

	want := &InsertStatement{
		Table: tokens[2],
		Values: &[]*Expression{
			{Literal: tokens[5], Kind: literalExpression},
			{Literal: tokens[7], Kind: literalExpression},
		},
	}
	got, newCursor, ok := parseInsertStatement(tokens, 0)
	assert.True(t, ok)
	assert.Equal(t, uint(9), newCursor)
	assert.Equal(t, got, want)
}

// CREATE TABLE users (id INT, name TEXT);
func TestCreateQuery(t *testing.T) {
	tokens := []*lexer.Token{
		{Value: "create", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 0}},
		{Value: "table", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 7}},
		{Value: "users", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 1, Col: 0}},
		{Value: "(", Type: lexer.SymbolType, Loc: lexer.Location{Line: 1, Col: 10}},
		{Value: "id", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 2, Col: 0}},
		{Value: "int", Type: lexer.KeywordType, Loc: lexer.Location{Line: 3, Col: 0}},
		{Value: ",", Type: lexer.SymbolType, Loc: lexer.Location{Line: 4, Col: 0}},
		{Value: "name", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 5, Col: 0}},
		{Value: "text", Type: lexer.KeywordType, Loc: lexer.Location{Line: 6, Col: 0}},
		{Value: ")", Type: lexer.SymbolType, Loc: lexer.Location{Line: 7, Col: 0}},
		{Value: ";", Type: lexer.SymbolType, Loc: lexer.Location{Line: 8, Col: 0}},
	}

	want := &CreateStatement{
		Table: tokens[2],
		ColumnDefs: &[]*ColumnDefinition{
			{ColName: *tokens[4], ColType: *tokens[5]},
			{ColName: *tokens[7], ColType: *tokens[8]},
		},
	}
	got, newCursor, ok := parseCreateStatement(tokens, 0)
	assert.True(t, ok)
	assert.Equal(t, uint(10), newCursor)
	assert.Equal(t, got, want)
}

func TestStatementParser(t *testing.T) {
	tokens := []*lexer.Token{
		{Value: "create", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 0}},
		{Value: "table", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 7}},
		{Value: "users", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 1, Col: 0}},
		{Value: "(", Type: lexer.SymbolType, Loc: lexer.Location{Line: 1, Col: 10}},
		{Value: "id", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 2, Col: 0}},
		{Value: "int", Type: lexer.KeywordType, Loc: lexer.Location{Line: 3, Col: 0}},
		{Value: ",", Type: lexer.SymbolType, Loc: lexer.Location{Line: 4, Col: 0}},
		{Value: "name", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 5, Col: 0}},
		{Value: "text", Type: lexer.KeywordType, Loc: lexer.Location{Line: 6, Col: 0}},
		{Value: ")", Type: lexer.SymbolType, Loc: lexer.Location{Line: 7, Col: 0}},
		{Value: ";", Type: lexer.SymbolType, Loc: lexer.Location{Line: 8, Col: 0}},
	}

	createStatement := &CreateStatement{
		Table: tokens[2],
		ColumnDefs: &[]*ColumnDefinition{
			{ColName: *tokens[4], ColType: *tokens[5]},
			{ColName: *tokens[7], ColType: *tokens[8]},
		},
	}
	want := &Statement{
		CreateStatement: createStatement,
		Kind:            CreateStateKind,
	}
	got, newCursor, ok := parseStatement(tokens, 0)
	assert.True(t, ok)
	assert.Equal(t, uint(10), newCursor)
	assert.Equal(t, got, want)
}

// CREATE TABLE users (id INT, name TEXT);
// INSERT INTO tablename VALUES (1, 2);
// SELECT id FROM users;
func TestParser(t *testing.T) {
	tokens := []*lexer.Token{
		// create: 0-10
		{Value: "create", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 0}},
		{Value: "table", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 7}},
		{Value: "users", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 1, Col: 0}},
		{Value: "(", Type: lexer.SymbolType, Loc: lexer.Location{Line: 1, Col: 10}},
		{Value: "id", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 2, Col: 0}},
		{Value: "int", Type: lexer.KeywordType, Loc: lexer.Location{Line: 3, Col: 0}},
		{Value: ",", Type: lexer.SymbolType, Loc: lexer.Location{Line: 4, Col: 0}},
		{Value: "name", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 5, Col: 0}},
		{Value: "text", Type: lexer.KeywordType, Loc: lexer.Location{Line: 6, Col: 0}},
		{Value: ")", Type: lexer.SymbolType, Loc: lexer.Location{Line: 7, Col: 0}},
		{Value: ";", Type: lexer.SymbolType, Loc: lexer.Location{Line: 8, Col: 0}},
		// insert: 11-20
		{Value: "insert", Type: lexer.KeywordType, Loc: lexer.Location{Line: 9, Col: 0}},
		{Value: "into", Type: lexer.KeywordType, Loc: lexer.Location{Line: 9, Col: 7}},
		{Value: "tablename", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 10, Col: 0}},
		{Value: "values", Type: lexer.KeywordType, Loc: lexer.Location{Line: 10, Col: 10}},
		{Value: "(", Type: lexer.SymbolType, Loc: lexer.Location{Line: 11, Col: 0}},
		{Value: "1", Type: lexer.NumType, Loc: lexer.Location{Line: 12, Col: 0}},
		{Value: ",", Type: lexer.SymbolType, Loc: lexer.Location{Line: 13, Col: 0}},
		{Value: "2", Type: lexer.NumType, Loc: lexer.Location{Line: 14, Col: 0}},
		{Value: ")", Type: lexer.SymbolType, Loc: lexer.Location{Line: 15, Col: 0}},
		{Value: ";", Type: lexer.SymbolType, Loc: lexer.Location{Line: 16, Col: 0}},
		// select: 21-25
		{Value: "select", Type: lexer.KeywordType, Loc: lexer.Location{Line: 17, Col: 0}},
		{Value: "id", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 17, Col: 7}},
		{Value: "from", Type: lexer.KeywordType, Loc: lexer.Location{Line: 18, Col: 0}},
		{Value: "users", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 18, Col: 5}},
		{Value: ";", Type: lexer.SymbolType, Loc: lexer.Location{Line: 18, Col: 10}},
	}

	want := &Ast{
		Statements: []*Statement{
			{
				CreateStatement: &CreateStatement{
					Table: tokens[2],
					ColumnDefs: &[]*ColumnDefinition{
						{ColName: *tokens[4], ColType: *tokens[5]},
						{ColName: *tokens[7], ColType: *tokens[8]},
					},
				},
				Kind: CreateStateKind,
			},
			{
				InsertStatement: &InsertStatement{
					Table: tokens[13],
					Values: &[]*Expression{
						{Literal: tokens[16], Kind: literalExpression},
						{Literal: tokens[18], Kind: literalExpression},
					},
				},
				Kind: InsertStateKind,
			},
			{
				SelectStatement: &SelectStatement{
					Items: &[]*Expression{
						{Literal: tokens[22], Kind: literalExpression},
					},
					From: tokens[24],
				},
				Kind: SelectStateKind,
			},
		},
	}

	got, err := Parse(tokens, 0)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
