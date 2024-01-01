package token

type TokenType string

type Token struct {
	Type 	TokenType
	Literal	string
}

const (
	LBRACE = "{"
	RBRACE = "}"
	LBRACK = "["
	RBRACK = "]"

	COMMA = ","
	DD = ":"
	QUOTES = "\""
	
	KEY = "KEY" // Any key in a key:value pair
	INT = "INT" // 1, 2, 34...
	STRING = "STRING" // Any value in a key:value pair
)
