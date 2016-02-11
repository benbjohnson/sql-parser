package main

import "strings"

func main() {
	query := "select id, name from employees where department_id = 3 and salary > 10"
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
	for _, condition := range stmt.Conditions {
		println(condition)
	}
}
