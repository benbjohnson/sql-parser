package sql

import (
	"strings"
	"testing"

	log "github.com/cihub/seelog"
	T "github.com/oldenbur/sql-parser/testutil"
	. "github.com/smartystreets/goconvey/convey"
)

func init() { T.ConfigureTestLogger() }

func TestParser_ParseSimpleSelect(t *testing.T) {

	Convey("Fields String method\n", t, func() {

		f := Fields{Field{Name: "f1"}}
		So(f.String(), ShouldEqual, "f1")

		f = Fields{Field{Name: "f1", Alias: "a1"}}
		So(f.String(), ShouldEqual, "f1 a1")

		f = Fields{Field{Name: "f1", Alias: "a1"}, Field{Name: "f2", Alias: "a2"}}
		So(f.String(), ShouldEqual, "f1 a1, f2 a2")

		f = Fields{Field{Name: "f1", Alias: "a1"}, Field{Name: "f2"}, Field{Name: "f3", Alias: "a3"}}
		So(f.String(), ShouldEqual, "f1 a1, f2, f3 a3")
	})

	Convey("Single field statement\n", t, func() {
		stmt, err := testParse(`SELECT name FROM tbl`)
		So(err, ShouldBeNil)
		log.Debug("SQL: ", stmt)
		So(stmt, ShouldResemble, &SelectStatement{
			FieldList: Fields{Field{Name: "name"}},
			TableList: Fields{Field{Name: "tbl"}},
		})
	})

	Convey("Multi-field statement\n", t, func() {
		stmt, err := testParse(`SELECT first_name, last_name, age FROM my_table`)
		So(err, ShouldBeNil)
		log.Debug("SQL: ", stmt)
		So(stmt, ShouldResemble, &SelectStatement{
			FieldList: Fields{Field{Name: "first_name"}, Field{Name: "last_name"}, Field{Name: "age"}},
			TableList: Fields{Field{Name: "my_table"}},
		})
	})

	Convey("Select all statement\n", t, func() {
		stmt, err := testParse(`SELECT * FROM my_table`)
		So(err, ShouldBeNil)
		log.Debug("SQL: ", stmt)
		So(stmt, ShouldResemble, &SelectStatement{
			FieldList: Fields{Field{Name: "*"}},
			TableList: Fields{Field{Name: "my_table"}},
		})
	})

	Convey("Expected SELECT", t, func() {
		_, err := testParse(`foo`)
		So(errstring(err), ShouldEqual, `found "foo", expected SELECT`)
	})

	Convey("Expected field", t, func() {
		_, err := testParse(`SELECT !`)
		So(errstring(err), ShouldEqual, `error parsing SELECT fields: found "!", expected field`)
	})

	Convey("Expected field", t, func() {
		_, err := testParse(`SELECT field1 alias1 BAD`)
		So(errstring(err), ShouldEqual, `found "BAD", expected FROM`)
	})

	Convey("Expected field", t, func() {
		f, err := testParse(`SELECT field1 alias1 FROM table1 talias1 BAD`)
		log.Debugf("f: %s", f)

		So(errstring(err), ShouldEqual, `found "BAD", expected WHERE`)
	})

}

func TestParseCommaDelimIdents(t *testing.T) {

	Convey("Expected table name", t, func() {
		p := NewParser(strings.NewReader(`field1`))
		f, err := p.parseCommaDelimIdents()
		So(err, ShouldBeNil)
		So(f, ShouldResemble, Fields{Field{Name: "field1", Alias: ""}})

		p = NewParser(strings.NewReader(`field1 alias1`))
		f, err = p.parseCommaDelimIdents()
		So(err, ShouldBeNil)
		So(f, ShouldResemble, Fields{Field{Name: "field1", Alias: "alias1"}})

		p = NewParser(strings.NewReader(`field1 alias1, field2 alias2`))
		f, err = p.parseCommaDelimIdents()
		So(err, ShouldBeNil)
		So(f, ShouldResemble, Fields{Field{Name: "field1", Alias: "alias1"},
			Field{Name: "field2", Alias: "alias2"}})

		p = NewParser(strings.NewReader(`field1 alias1, field2 alias2, field3`))
		f, err = p.parseCommaDelimIdents()
		So(err, ShouldBeNil)
		So(f, ShouldResemble, Fields{Field{Name: "field1", Alias: "alias1"},
			Field{Name: "field2", Alias: "alias2"},
			Field{Name: "field3"}})

		p = NewParser(strings.NewReader(`field1 alias1, field2`))
		f, err = p.parseCommaDelimIdents()
		So(err, ShouldBeNil)
		So(f, ShouldResemble, Fields{Field{Name: "field1", Alias: "alias1"},
			Field{Name: "field2"}})

	})

}


// testParse returns the result of parsing the specified string.
func testParse(s string) (*SelectStatement, error) {
	return NewParser(strings.NewReader(s)).Parse()
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
