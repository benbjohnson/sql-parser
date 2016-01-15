package sql

import (
	"fmt"
	"io"

//	log "github.com/cihub/seelog"
)

type Field struct {
	Name string
	Alias string
}

func (f Field) String() string {
	if len(f.Alias) < 1 {
		return f.Name
	} else {
		return fmt.Sprintf("%s %s", f.Name, f.Alias)
	}
}

type Fields []Field

func (f Fields) String() string {
	if len(f) < 1 {
		return ""
	} else if len(f) == 1 {
		return f[0].String()
	} else {
		return fmt.Sprintf("%s, %s", f[0], f[1:])
	}
}

// SelectStatement represents a SQL SELECT statement.
type SelectStatement struct {
	FieldList Fields
	TableList Fields
}

func (s SelectStatement) String() string {
	return fmt.Sprintf("SELECT %s FROM %s", s.FieldList.String(), s.TableList.String())
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a SQL SELECT statement.
func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}

	selFields, err := p.parseCommaDelimIdents()
	if err != nil {
		return nil, fmt.Errorf("error parsing SELECT fields: %v", err)
	}
	stmt.FieldList = selFields

	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}

	tables, err := p.parseCommaDelimIdents()
	if err != nil {
		return nil, fmt.Errorf("error parsing SELECT fields: %v", err)
	}
	stmt.TableList = tables

	// Next we should see the "HWERE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != WHERE && tok != EOF{
		return nil, fmt.Errorf("found %q, expected WHERE", lit)
	}

	// Return the successfully parsed statement.
	return stmt, nil
}

// parseCommaDelimIdents assumes that the scanner position is at the head
// of comma-delimited list of fields each possibly followed by an alias.
// The list is parsed int a Fields and returned along with any error that
// arises during parsing.
func (p *Parser) parseCommaDelimIdents() (fields Fields, err error) {

	for {

		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}

		f := Field{ Name: lit }

		tok, lit = p.scanIgnoreWhitespace()
		if tok == IDENT {
			f.Alias = lit
			tok, lit = p.scanIgnoreWhitespace()
		}

		fields = append(fields, f)

		if tok != COMMA {
			p.unscan()
			break
		}
	}

	return
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
