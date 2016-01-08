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

	log.Debug("TestScanner")

	Convey("Special tokens\n", t, func() {
		testScanner(``, EOF, ``)
		testScanner(`#`, ILLEGAL, `#`)
		testScanner(` `, WS, ` `)
		testScanner(`    `, WS, `    `)
		testScanner("\t", WS, "\t")
		testScanner("\n", WS, "\n")
	})

	//	var tests = []struct {
	//		s   string
	//		tok sql.Token
	//		lit string
	//	}{
	//		// Special tokens (EOF, ILLEGAL, WS)
	//		{s: ``, tok: sql.EOF},
	//		{s: `#`, tok: sql.ILLEGAL, lit: `#`},
	//		{s: ` `, tok: sql.WS, lit: " "},
	//		{s: "\t", tok: sql.WS, lit: "\t"},
	//		{s: "\n", tok: sql.WS, lit: "\n"},
	//
	//		// Misc characters
	//		{s: `*`, tok: sql.ASTERISK, lit: "*"},
	//
	//		// Identifiers
	//		{s: `foo`, tok: sql.IDENT, lit: `foo`},
	//		{s: `Zx12_3U_-`, tok: sql.IDENT, lit: `Zx12_3U_`},
	//
	//		// Keywords
	//		{s: `FROM`, tok: sql.FROM, lit: "FROM"},
	//		{s: `SELECT`, tok: sql.SELECT, lit: "SELECT"},
	//	}
	//
	//	for i, tt := range tests {
	//		s := sql.NewScanner(strings.NewReader(tt.s))
	//		tok, lit := s.Scan()
	//		if tt.tok != tok {
	//			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
	//		} else if tt.lit != lit {
	//			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
	//		}
	//	}
}

func testScanner(str string, tok Token, lit string) {
	s := NewScanner(strings.NewReader(str))
	tokTest, litTest := s.Scan()
	So(tokTest, ShouldEqual, tok)
	So(litTest, ShouldEqual, lit)
}
