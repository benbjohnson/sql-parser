package sql

import (
	"fmt"
)

type Cond interface{}

// CondComp represents a single comparison, e.g. f = 'bucky'
type CondComp struct {
	Ident string
	CondOp Token  // e.g. =, <=
	Val string
}

func (c CondComp) String() string {
	return fmt.Sprintf("%s", c.Ident)
//	return fmt.Sprintf("Comp{%s %s %s}", c.Ident, c.CondOp, c.Val)
}

// CondConj represents a single level of ANDed or ORed statements,
// e.g. f1 = "v1" AND myNum >= 12.34 AND (f2 != "v2" OR id = 12)
// There is an AND node with two Conds and a single Node, which is
// itself an OR node with two conds and no Nodes.
type CondConj struct {
	Op Token  // AND or OR
	Left Cond
	Right Cond
}

func (c CondConj) String() string {
	return fmt.Sprintf("(%s %s %s)", c.Left, c.Op, c.Right)
}

func (p *Parser) parseCondTree() (Cond, error) {

	var left, right Cond
	var err error

	tok, lit := p.scanIgnoreWhitespace()
	if tok == PAREN_L {
		left, err = p.parseCondTree()
		if err != nil {
			return nil, err
		}
	} else if tok == IDENT {
		left = &CondComp{Ident: lit}
	} else {
		return nil, fmt.Errorf(`expeected PAREN_L or IDENT, got "%s"`, lit)
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok == AND || tok == OR {

		condConj := &CondConj{Op: tok, Left: left}
		right, err = p.parseCondTree()
		if err != nil {
			return nil, err
		}
		condConj.Right = right

		if tok, lit = p.scanIgnoreWhitespace(); tok != PAREN_R {
			p.unscan()
		}

		return condConj, nil
	} else if tok == PAREN_R {
		return left, nil
	} else if tok != EOF {
		return nil, fmt.Errorf(`expected AND or OR, got "%s"`, lit)
	}

	return left, nil
}