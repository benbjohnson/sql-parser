package sql

import (
	"strings"
	"testing"

	log "github.com/cihub/seelog"
	T "github.com/oldenbur/sql-parser/testutil"
	. "github.com/smartystreets/goconvey/convey"
)

func init() { T.ConfigureTestLogger() }

func TestParserCond(t *testing.T) {

	defer log.Flush()

	Convey("Test parsing basic conditions\n", t, func() {
		p := NewParser(strings.NewReader(`A`))
		c, err := p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondComp{Ident:"A"})
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`A AND B`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondConj{Left: &CondComp{Ident:"A"}, Op: AND, Right: &CondComp{Ident:"B"}})
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`A AND B AND C AND D`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondConj{Left: &CondComp{Ident:"A"}, Op: AND,
			Right: &CondConj{Left: &CondComp{Ident:"B"}, Op: AND,
				Right: &CondConj{Left: &CondComp{Ident:"C"}, Op: AND,
					Right: &CondComp{Ident:"D"}}}})
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`(A AND B)`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondConj{Left: &CondComp{Ident:"A"}, Op: AND, Right: &CondComp{Ident:"B"}})
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`(A AND B) OR C`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondConj{
			Left: &CondConj{Left: &CondComp{Ident:"A"}, Op: AND, Right: &CondComp{Ident:"B"}},
			Op: OR,
			Right: &CondComp{Ident:"C"},
		})
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`(A AND B) OR (C AND D)`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondConj{
			Left: &CondConj{Left: &CondComp{Ident:"A"}, Op: AND, Right: &CondComp{Ident:"B"}},
			Op: OR,
			Right: &CondConj{Left: &CondComp{Ident:"C"}, Op: AND, Right: &CondComp{Ident:"D"}},
		})
		log.Debugf("cond: %s", c)

		p = NewParser(strings.NewReader(`(A AND B) OR (C AND (D OR E) AND F)`))
		c, err = p.parseCondTree()
		So(err, ShouldBeNil)
		So(c, ShouldResemble, &CondConj{
			Left: &CondConj{Left: &CondComp{Ident:"A"}, Op: AND, Right: &CondComp{Ident:"B"}},
			Op: OR,
			Right: &CondConj{Left: &CondComp{Ident:"C"}, Op: AND,
					Right: &CondConj{Left: &CondConj{Left: &CondComp{Ident:"D"}, Op: OR, Right: &CondComp{Ident:"E"}}, Op: AND,
									Right: &CondComp{Ident:"F"}}}})
		log.Debugf("cond: %s", c)

	})



}