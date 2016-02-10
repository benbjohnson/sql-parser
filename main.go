package main

import "strings"

func main() {
	query := "select * from my_table"
	println(query)
	r := strings.NewReader(query)

	parser := NewParser(r)

	stmt, err := parser.Parse()
	if err != nil {
		println(err.Error())
	}

	println(stmt.TableName)
	for _, field := range stmt.Fields {
		println(field)
	}
}
