package backend

import (
	"databaseInGo/lexer"
	"databaseInGo/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) *Mb {
	t.Helper() // Marks this as a helper so failure line numbers point to the test
	db := NewMb()
	return db
}

// CREATE TABLE users (id INT, name TEXT);
func usersCreateStatement() *parser.CreateStatement {
	return &parser.CreateStatement{
		Table: &lexer.Token{Value: "users", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 0, Col: 13}},
		ColumnDefs: &[]*parser.ColumnDefinition{
			{
				ColName: lexer.Token{Value: "id", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 0, Col: 20}},
				ColType: lexer.Token{Value: "int", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 23}},
			},
			{
				ColName: lexer.Token{Value: "name", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 0, Col: 28}},
				ColType: lexer.Token{Value: "text", Type: lexer.KeywordType, Loc: lexer.Location{Line: 0, Col: 33}},
			},
		},
	}
}

// INSERT INTO users VALUES (1, 'alice');
func usersInsertStatement() *parser.InsertStatement {
	return &parser.InsertStatement{
		Table: &lexer.Token{Value: "users", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 0, Col: 12}},
		Values: &[]*parser.Expression{
			{Literal: &lexer.Token{Value: "1", Type: lexer.NumType, Loc: lexer.Location{Line: 0, Col: 26}}},
			{Literal: &lexer.Token{Value: "alice", Type: lexer.TextType, Loc: lexer.Location{Line: 0, Col: 29}}},
		},
	}
}

func TestValidExecCreate(t *testing.T) {
	db := setup(t)

	want := &Table{
		columns: []*Column{
			{ColName: "id", ColType: "int"},
			{ColName: "name", ColType: "text"},
		},
	}

	err := db.ExecCreate(usersCreateStatement())
	assert.NoError(t, err)
	assert.Equal(t, want, db.Tables["users"])
}

func TestValidExecInsert(t *testing.T) {
	db := setup(t)
	assert.NoError(t, db.ExecCreate(usersCreateStatement()))

	want := &Table{
		columns: []*Column{
			{ColName: "id", ColType: "int"},
			{ColName: "name", ColType: "text"},
		},
		rows: [][]memCell{
			{memCell{0, 0, 0, 1}, memCell("alice")},
		},
	}

	err := db.ExecInsert(usersInsertStatement())
	assert.NoError(t, err)
	assert.Equal(t, want, db.Tables["users"])
}

// SELECT id, name FROM users;
func TestValidExecSelect(t *testing.T) {
	db := setup(t)
	assert.NoError(t, db.ExecCreate(usersCreateStatement()))
	assert.NoError(t, db.ExecInsert(usersInsertStatement()))

	ss := &parser.SelectStatement{
		Items: &[]*parser.Expression{
			{Literal: &lexer.Token{Value: "id", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 0, Col: 7}}},
			{Literal: &lexer.Token{Value: "name", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 0, Col: 11}}},
		},
		From: &lexer.Token{Value: "users", Type: lexer.IdentifierType, Loc: lexer.Location{Line: 0, Col: 21}},
	}

	want := &Result{
		Columns: []*Column{
			{ColName: "id", ColType: "int"},
			{ColName: "name", ColType: "text"},
		},
		Rows: [][]Cell{
			{memCell{0, 0, 0, 1}, memCell("alice")},
		},
	}

	got, err := db.ExecSelect(ss)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
