package wal

import (
	"bytes"
	"databaseInGo/backend"
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeEncode(t *testing.T) {
	stmt := "INSERT INTO users VALUES (1, 'alice');"
	encStmt := EncodeStatement(stmt)
	assert.Equal(t, len(encStmt), len(stmt)+8)
	dec := NewDecoder(bytes.NewReader(encStmt))
	rec, err := dec.Next()
	if err != nil {
		t.Fatalf("Next() error = %v, wantd nil", err)
	}
	if rec.data != stmt {
		t.Errorf("data = %q, want %q", rec.data, stmt)
	}
	if rec.dataLength != uint32(len(stmt)) {
		t.Errorf("dataLength = %d, want %d", rec.dataLength, len(stmt))
	}
	if _, err := dec.Next(); !errors.Is(err, io.EOF) {
		t.Errorf("second Next() error = %v, want io.EOF", err)
	}
}

func TestReplay(t *testing.T) {
	stmts := []string{"CREATE TABLE users (id INT, name TEXT);", "INSERT INTO users VALUES (1, 'phil');"}
	mb := backend.NewMb()
	for _, stmt := range stmts {
		mb.Execute(stmt)
	}

	walMb := backend.NewMb()
	testDir := t.TempDir()
	walPath := filepath.Join(testDir, "testWal.log")
	wal, err := NewWal(walPath)
	if err != nil {
		t.Fatal("error creating new wal")
	}
	for _, stmt := range stmts {
		err = wal.Append(stmt)
		if err != nil {
			t.Errorf("expected error to be not nil, actual %s", err.Error())
		}
	}
	err = wal.Replay(walMb)
	if err != nil {
		t.Errorf("expected error to be not nil, actual %s", err.Error())
	}
	assert.Equal(t, mb, walMb)
}
