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
	NUMBER // e.g. 1, 12.34, -46, -98.765

	// Misc characters
	ASTERISK // *
	COMMA    // ,
	PAREN_L  // (
	PAREN_R  // )

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
