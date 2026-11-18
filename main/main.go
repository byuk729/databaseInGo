package main

import (
	"bufio"
	"databaseInGo/backend"
	"databaseInGo/lexer"
	"databaseInGo/wal"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func printResult(r *backend.Result) {
	data := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	columnRow := []string{}
	for _, col := range r.Columns {
		columnRow = append(columnRow, col.ColName)
	}
	data = append(data, columnRow)
	table.Header(data[0])
	for _, row := range r.Rows {
		dataRow := []string{}
		for c := 0; c < len(row); c++ {
			cellType := r.Columns[c].ColType
			switch cellType {
			case string(lexer.IntKeyword):
				s := fmt.Sprintf("%d", row[c].AsInt())
				dataRow = append(dataRow, s)
			case string(lexer.TextKeyword):
				s := row[c].AsText()
				dataRow = append(dataRow, s)
			}
		}
		data = append(data, dataRow)
	}
	table.Bulk(data[1:])
	table.Render()
}

func printPrompt() {
	fmt.Print("Enter database query: ")
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func cleanInput(s string) string {
	res := strings.TrimSpace(s)
	return res
}

func main() {
	printPrompt()
	reader := bufio.NewScanner(os.Stdin)
	mb := backend.NewMb()
	wrtal, err := wal.NewWal("wal.log")
	if err != nil {
		panic(err)
	}
	defer wrtal.File.Close()
	err = wrtal.Replay(mb)
	if err != nil {
		panic(err)
	}
	for reader.Scan() {
		text := cleanInput(reader.Text())

		ast, err := backend.Parse(text)
		if err != nil {
			fmt.Println("error:", err)
			printPrompt()
			continue
		}

		res, err := mb.ExecuteAst(ast)
		if err != nil {
			fmt.Println("Unable to execute query:", err)
			printPrompt()
			continue
		}

		if ast.HasWrites() {
			if err := wrtal.Append(text); err != nil {
				panic(err)
			}
			fmt.Println("ok")
		}

		if res != nil {
			printResult(res)
			fmt.Println()
		}
		printPrompt()
	}
	if err := reader.Err(); err != nil {
		fmt.Print("there was an error")
	}
}
