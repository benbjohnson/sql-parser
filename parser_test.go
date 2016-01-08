package sql_test

import (
	"strings"
	"testing"

	log "github.com/cihub/seelog"
	"github.com/oldenbur/sql-parser"
	T "github.com/oldenbur/sql-parser/testutil"
	. "github.com/smartystreets/goconvey/convey"
)

func init() { T.ConfigureTestLogger() }

func TestParser_ParseSimpleSelect(t *testing.T) {

	Convey("Single field statement\n", t, func() {
		stmt, err := testParse(`SELECT name FROM tbl`)
		So(err, ShouldBeNil)
		log.Debug("SQL: ", stmt)
		So(stmt, ShouldResemble, &sql.SelectStatement{Fields: []string{"name"}, TableName: "tbl"})
	})

	Convey("Multi-field statement\n", t, func() {
		stmt, err := testParse(`SELECT first_name, last_name, age FROM my_table`)
		So(err, ShouldBeNil)
		log.Debug("SQL: ", stmt)
		So(stmt, ShouldResemble, &sql.SelectStatement{
			Fields:    []string{"first_name", "last_name", "age"},
			TableName: "my_table",
		})
	})

	Convey("Select all statement\n", t, func() {
		stmt, err := testParse(`SELECT * FROM my_table`)
		So(err, ShouldBeNil)
		log.Debug("SQL: ", stmt)
		So(stmt, ShouldResemble, &sql.SelectStatement{Fields: []string{"*"}, TableName: "my_table"})
	})

	Convey("Expected SELECT", t, func() {
		_, err := testParse(`foo`)
		So(errstring(err), ShouldEqual, `found "foo", expected SELECT`)
	})

	Convey("Expected field", t, func() {
		_, err := testParse(`SELECT !`)
		So(errstring(err), ShouldEqual, `found "!", expected field`)
	})

	Convey("Expected FROM", t, func() {
		_, err := testParse(`SELECT field xxx`)
		So(errstring(err), ShouldEqual, `found "xxx", expected FROM`)
	})

	Convey("Expected table name", t, func() {
		_, err := testParse(`SELECT field FROM *`)
		So(errstring(err), ShouldEqual, `found "*", expected table name`)
	})
}

// testParse returns the result of parsing the specified string.
func testParse(s string) (*sql.SelectStatement, error) {
	return sql.NewParser(strings.NewReader(s)).Parse()
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
