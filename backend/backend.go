package backend

import (
	"bytes"
	"databaseInGo/lexer"
	"databaseInGo/parser"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

var (
	TableAlreadyExists      error = errors.New("Table already exists")
	TableDoesNotExist       error = errors.New("Table does not exist")
	mismatchedValues        error = errors.New("Table columns and statement values do not match in size")
	missingValuesForInsert  error = errors.New("Missing values for insert")
	columnValueTypeMismatch error = errors.New("Column type doesn't match value type")
	columnNotFound          error = errors.New("select item has no matching column")
)

type Result struct {
	Columns []*Column
	Rows    [][]Cell
}

type Column struct {
	ColName string
	ColType string
}

type Cell interface {
	AsInt() int32
	AsText() string
}

type dbInterface interface {
	ExecSelect(*parser.SelectStatement) (*Result, error)
	ExecInsert(*parser.InsertStatement) error
	ExecCreate(*parser.CreateStatement) error
}

type memCell []byte

func (mc memCell) AsInt() int32 {
	var i int32
	buf := bytes.NewReader(mc)
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		fmt.Println("binary.Read for AsInt() failed:", err)
	}
	return i
}

func (mc memCell) AsText() string {
	return string(mc)
}

type Table struct {
	columns []*Column
	rows    [][]memCell
}

type Mb struct {
	Tables map[string]*Table
}

func NewMb() *Mb {
	return &Mb{
		Tables: make(map[string]*Table),
	}
}

func (db *Mb) ExecCreate(cs *parser.CreateStatement) error {
	if *cs.ColumnDefs == nil {
		return nil
	}
	TableName := cs.Table.Value
	if _, ok := db.Tables[TableName]; ok {
		return TableAlreadyExists
	}
	tbl := &Table{}
	cols := []*Column{}
	colDefs := cs.ColumnDefs
	for _, colDef := range *colDefs {
		col := &Column{}
		col.ColName = colDef.ColName.Value
		col.ColType = colDef.ColType.Value
		cols = append(cols, col)
	}
	tbl.columns = cols
	db.Tables[TableName] = tbl
	return nil
}

func toBytes(t *lexer.Token) (memCell, error) {
	if t.Type == lexer.TextType {
		return []byte(t.Value), nil
	} else if t.Type == lexer.NumType {
		i64, err := strconv.ParseInt(t.Value, 10, 32)
		if err != nil {
			return nil, err
		}
		i := int32(i64)
		buf := bytes.NewBuffer(make([]byte, 0, 4))
		err = binary.Write(buf, binary.BigEndian, i)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return []byte{}, errors.New("Invalid token type")
}

func matchingTypes(columnType string, valueType lexer.TokenType) bool {
	switch lexer.Keyword(columnType) {
	case lexer.IntKeyword:
		return valueType == lexer.NumType
	case lexer.TextKeyword:
		return valueType == lexer.TextType
	default:
		return false
	}
}

func (db *Mb) ExecInsert(is *parser.InsertStatement) error {
	if is.Values == nil || len(*is.Values) == 0 {
		return missingValuesForInsert
	}
	TableName := is.Table.Value
	if _, ok := db.Tables[TableName]; !ok {
		return TableDoesNotExist
	}
	curTable := db.Tables[TableName]
	if len(curTable.columns) != len(*is.Values) {
		return mismatchedValues
	}
	row := []memCell{}
	for i := 0; i < len(curTable.columns); i++ {
		columnType := curTable.columns[i].ColType
		valueType := (*is.Values)[i].Literal.Type
		if !matchingTypes(columnType, valueType) {
			return columnValueTypeMismatch
		}
		curValue := (*is.Values)[i].Literal
		curByte, err := toBytes(curValue)
		if err != nil {
			return err
		}
		row = append(row, curByte)
	}
	curTable.rows = append(curTable.rows, row)
	return nil
}

func (db *Mb) ExecSelect(ss *parser.SelectStatement) (*Result, error) {
	if ss.Items == nil {
		return nil, nil
	}
	TableName := ss.From.Value
	if _, ok := db.Tables[TableName]; !ok {
		return nil, TableDoesNotExist
	}
	curTable := db.Tables[TableName]
	res := Result{}

	// Find each selected column once, so this works even when the Table has no rows.
	colIdxs := []int{}
items:
	for _, exps := range *ss.Items {
		lCol := exps.Literal.Value
		for j, col := range curTable.columns {
			if lCol == col.ColName {
				res.Columns = append(res.Columns, &Column{col.ColName, col.ColType})
				colIdxs = append(colIdxs, j)
				continue items
			}
		}
		return nil, columnNotFound
	}

	for _, row := range curTable.rows {
		resRow := []Cell{}
		for _, j := range colIdxs {
			resRow = append(resRow, row[j])
		}
		res.Rows = append(res.Rows, resRow)
	}
	return &res, nil
}

func Parse(text string) (*parser.Ast, error) {
	tokens, err := lexer.Lex(text)
	if err != nil {
		return nil, err
	}
	return parser.Parse(tokens, 0)
}

func (db *Mb) ExecuteAst(ast *parser.Ast) (*Result, error) {
	var res *Result
	for _, statement := range ast.Statements {
		switch statement.Kind {
		case parser.CreateStateKind:
			if err := db.ExecCreate(statement.CreateStatement); err != nil {
				return nil, err
			}
		case parser.InsertStateKind:
			if err := db.ExecInsert(statement.InsertStatement); err != nil {
				return nil, err
			}
		case parser.SelectStateKind:
			r, err := db.ExecSelect(statement.SelectStatement)
			if err != nil {
				return nil, err
			}
			res = r
		}
	}
	return res, nil
}

func (db *Mb) Execute(text string) (*Result, error) {
	ast, err := Parse(text)
	if err != nil {
		return nil, err
	}
	return db.ExecuteAst(ast)
}
