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
		testScanString(``, EOF, ``)
		testScanString(`#`, ILLEGAL, `#`)
		testScanString(` `, WS, ` `)
		testScanString(`    `, WS, `    `)
		testScanString("\t", WS, "\t")
		testScanString("\n", WS, "\n")

		testScanString(`*`, ASTERISK, `*`)
		testScanString(`,`, COMMA, `,`)
		testScanString(`(`, PAREN_L, `(`)
		testScanString(`)`, PAREN_R, `)`)
	})

	Convey("Identifiers\n", t, func() {
		testScanString(`foo`, IDENT, `foo`)
		testScanString(`Zx12_3U_-*`, IDENT, `Zx12_3U_-*`)
	})

	Convey("Numbers\n", t, func() {
		testScanString(`1`, NUMBER, `1`)
		testScanString(`12.34`, NUMBER, `12.34`)
		testScanString(`-46`, NUMBER, `-46`)
		testScanString(`-98.765`, NUMBER, `-98.765`)
	})

	Convey("Keywords\n", t, func() {
		testScanString(`FROM`, FROM, `FROM`)
		testScanString(`From`, FROM, `From`)
		testScanString(`from`, FROM, `from`)
		testScanString(`SELECT`, SELECT, `SELECT`)
		testScanString(`WHERE`, WHERE, `WHERE`)
		testScanString(`AND`, AND, `AND`)
		testScanString(`OR`, OR, `OR`)
	})

	Convey("Operators\n", t, func() {
		testScanString(`=`, EQ, `=`)
		testScanString(`!=`, NE, `!=`)
		testScanString(`<`, LT, `<`)
		testScanString(`>`, GT, `>`)
		testScanString(`<=`, LE, `<=`)
		testScanString(`>=`, GE, `>=`)
	})

	Convey("Strings\n", t, func() {
		testScanString(`"abc"`, STRING, `"abc"`)
		testScanString(`'abc 789 @#$'`, STRING, `'abc 789 @#$'`)
		testScanString(`"Howdy\" \"ho"`, STRING, `"Howdy\" \"ho"`)
		testScanString(`'""Bucky  \'"\' Badger""'`, STRING, `'""Bucky  \'"\' Badger""'`)
		testScanString(`'illegal1
		'`, ILLEGAL, `'illegal1`)
		testScanString(`'illegal2`, ILLEGAL, `'illegal2`)
	})

	Convey("Real statement - somewhat complicated\n", t, func() {
		str := `SELECT t1.field1, t2.* FROM table1 t1
				wHeRe t1.joinA = t2.joinA AND (t2.fieldN <= -123.456 OR t2.fieldS = 'howdy ho')`
		s := NewScanner(strings.NewReader(str))

		testScanRmWs(s, SELECT, `SELECT`)
		testScanRmWs(s, IDENT, `t1.field1`)
		testScanRmWs(s, COMMA, `,`)
		testScanRmWs(s, IDENT, `t2.*`)
		testScanRmWs(s, FROM, `FROM`)
		testScanRmWs(s, IDENT, `table1`)
		testScanRmWs(s, IDENT, `t1`)
		testScanRmWs(s, WHERE, `wHeRe`)
		testScanRmWs(s, IDENT, `t1.joinA`)
		testScanRmWs(s, EQ, `=`)
		testScanRmWs(s, IDENT, `t2.joinA`)
		testScanRmWs(s, AND, `AND`)
		testScanRmWs(s, PAREN_L, `(`)
		testScanRmWs(s, IDENT, `t2.fieldN`)
		testScanRmWs(s, LE, `<=`)
		testScanRmWs(s, NUMBER, `-123.456`)
		testScanRmWs(s, OR, `OR`)
		testScanRmWs(s, IDENT, `t2.fieldS`)
		testScanRmWs(s, EQ, `=`)
		testScanRmWs(s, STRING, `'howdy ho'`)
		testScanRmWs(s, PAREN_R, `)`)
	})
}

func testScanString(str string, tok Token, lit string) {
	s := NewScanner(strings.NewReader(str))
	tokTest, litTest := s.Scan()
	So(tokTest, ShouldEqual, tok)
	So(litTest, ShouldEqual, lit)
}

func testScanRmWs(s *Scanner, tok Token, lit string) {
	litTest := ""
	tokTest := WS
	for tokTest == WS {
		tokTest, litTest = s.Scan()
	}
	So(tokTest, ShouldEqual, tok)
	So(litTest, ShouldEqual, lit)
}
