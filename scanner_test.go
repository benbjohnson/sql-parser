package sql

import (
	"strings"
	"testing"

	log "github.com/cihub/seelog"
	T "github.com/oldenbur/sql-parser/testutil"
	. "github.com/smartystreets/goconvey/convey"
)

func init() { T.ConfigureTestLogger() }

func TestScanner(t *testing.T) {

	defer log.Flush()

	Convey("Special tokens\n", t, func() {
		testScanner(``, EOF, ``)
		testScanner(`#`, ILLEGAL, `#`)
		testScanner(` `, WS, ` `)
		testScanner(`    `, WS, `    `)
		testScanner("\t", WS, "\t")
		testScanner("\n", WS, "\n")

		testScanner(`*`, ASTERISK, `*`)
		testScanner(`,`, COMMA, `,`)
		testScanner(`(`, PAREN_L, `(`)
		testScanner(`)`, PAREN_R, `)`)
	})

	Convey("Identifiers\n", t, func() {
		testScanner(`foo`, IDENT, `foo`)
		testScanner(`Zx12_3U_-*`, IDENT, `Zx12_3U_-*`)
	})

	Convey("Numbers\n", t, func() {
		testScanner(`1`, NUMBER, `1`)
		testScanner(`12.34`, NUMBER, `12.34`)
		testScanner(`-46`, NUMBER, `-46`)
		testScanner(`-98.765`, NUMBER, `-98.765`)
	})

	Convey("Keywords\n", t, func() {
		testScanner(`FROM`, FROM, `FROM`)
		testScanner(`From`, FROM, `From`)
		testScanner(`from`, FROM, `from`)
		testScanner(`SELECT`, SELECT, `SELECT`)
		testScanner(`WHERE`, WHERE, `WHERE`)
		testScanner(`AND`, AND, `AND`)
		testScanner(`OR`, OR, `OR`)
	})

	Convey("Operators\n", t, func() {
		testScanner(`=`, EQ, `=`)
		testScanner(`!=`, NE, `!=`)
		testScanner(`<`, LT, `<`)
		testScanner(`>`, GT, `>`)
		testScanner(`<=`, LE, `<=`)
		testScanner(`>=`, GE, `>=`)
	})

	Convey("Strings\n", t, func() {
		testScanner(`"abc"`, STRING, `"abc"`)
		testScanner(`'abc 789 @#$'`, STRING, `'abc 789 @#$'`)
		testScanner(`"Howdy\" \"ho"`, STRING, `"Howdy\" \"ho"`)
		testScanner(`'""Bucky  \'"\' Badger""'`, STRING, `'""Bucky  \'"\' Badger""'`)
		testScanner(`'illegal1
		'`, ILLEGAL, `'illegal1`)
		testScanner(`'illegal2`, ILLEGAL, `'illegal2`)
	})

}

func testScanner(str string, tok Token, lit string) {
	s := NewScanner(strings.NewReader(str))
	tokTest, litTest := s.Scan()
	So(tokTest, ShouldEqual, tok)
	So(litTest, ShouldEqual, lit)
}
