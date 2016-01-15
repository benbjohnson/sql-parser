package sql

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT  // main
	NUMBER // 1, 12.34, -46, -98.765
    STRING // 'abc', "DEF 123 &*$"

	// Misc characters
	ASTERISK   // *
	COMMA      // ,
	PAREN_L    // (
	PAREN_R    // )

	// Operators
	EQ // =
	NE // !=
	GT // >
	LT // <
	GE // >=
	LE // <=

	// Keywords
	SELECT
	FROM
	WHERE
	AND
	OR
)

func (t Token) String() string {
	switch (t) {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case ASTERISK:
		return "ASTERISK"
	case COMMA:
		return "COMMA"
	case PAREN_L:
		return "PAREN_L"
	case PAREN_R:
		return "PAREN_R"
	case EQ:
		return "EQ"
	case NE:
		return "NE"
	case GT:
		return "GT"
	case LT:
		return "LT"
	case GE:
		return "GE"
	case LE:
		return "LE"
	case SELECT:
		return "SELECT"
	case FROM:
		return "FROM"
	case WHERE:
		return "WHERE"
	case AND:
		return "AND"
	case OR:
		return "OR"
	}
	return "UNKNOWN"
}