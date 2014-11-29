package sql

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // main

	// Misc characters
	ASTERISK // *
	COMMA    // ,

	// Keywords
	SELECT
	FROM
)
