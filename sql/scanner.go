package sql

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isDigit(ch) || ch == '-' {
		s.unread()
		return s.scanNumber()
	} else if isOpChar(ch) {
		s.unread()
		return s.scanOp()
	} else if ch == '\'' {
		s.unread()
		return s.scanStrSngl()
	} else if ch == '"' {
		s.unread()
		return s.scanStrDbl()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	case '(':
		return PAREN_L, string(ch)
	case ')':
		return PAREN_R, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isIdentChar(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "SELECT":
		return SELECT, buf.String()
	case "FROM":
		return FROM, buf.String()
	case "WHERE":
		return WHERE, buf.String()
	case "AND":
		return AND, buf.String()
	case "OR":
		return OR, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

// scanNumber consumes the current rune and all contiguous number runes.
func (s *Scanner) scanNumber() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) && ch != '-' && ch != '.' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return NUMBER, buf.String()
}

// scanOp consumes the current rune and subsequent operator runes, returning
// the operator type and literal string, i.e. =, !=, <, >, <= or >=.
func (s *Scanner) scanOp() (tok Token, lit string) {
	var buf bytes.Buffer
	ch := s.read()
	buf.WriteRune(ch)

	if ch == '!' {
		ch = s.read()
		if ch != '=' {
			if ch != eof {
				buf.WriteRune(ch)
			}
			return ILLEGAL, buf.String()
		}
		buf.WriteRune(ch)
		return NE, buf.String()
	} else if ch == '=' {
		return EQ, buf.String()
	} else if ch == '<' {
		ch = s.read()
		if ch != '=' {
			s.unread()
			return LT, buf.String()
		} else {
			buf.WriteRune(ch)
			return LE, buf.String()
		}
		return
	} else if ch == '>' {
		ch = s.read()
		if ch != '=' {
			s.unread()
			return GT, buf.String()
		} else {
			buf.WriteRune(ch)
			return GE, buf.String()
		}
	}

	return ILLEGAL, buf.String()
}

// scanStr consumes the current rune, which is assumed to be a quote,
// and continues to consume until either a newline or an unescaped
// closing quote is encountered.
func (s *Scanner) scanStrSngl() (tok Token, lit string) {
	return s.scanStr('\'')
}

func (s *Scanner) scanStrDbl() (tok Token, lit string) {
	return s.scanStr('"')
}

func (s *Scanner) scanStr(term rune) (tok Token, lit string) {

	var buf bytes.Buffer
	ch := s.read()
	buf.WriteRune(ch)

	for {
		if ch := s.read(); ch == eof || ch == '\n' {
			return ILLEGAL, buf.String()
		} else if ch == '\\' {
			ch := s.read()
			buf.WriteRune('\\')
			buf.WriteRune(ch)
		} else if ch == term {
			buf.WriteRune(ch)
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return STRING, buf.String()
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) peek() rune {
	ch := s.read()
	s.unread()
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

// isIdentChar returns true if the run is a valid identifier character.
func isIdentChar(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_' || ch == '.' || ch == '-' || ch == '*'
}

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// isOpChar returns true if the run is an operator character.
func isOpChar(ch rune) bool { return ch == '=' || ch == '!' || ch == '<' || ch == '>' }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
