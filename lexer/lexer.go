package lexer

import (
	"github.com/cs50-romain/jsontogo/token"
)

type Lexer struct {
	input		string
	position	int
	readPosition	int
	ch		byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	if l.ch == '{' {
		tok = newToken(token.LBRACE, l.ch)
	} else if l.ch == '}' {
		tok = newToken(token.RBRACE, l.ch)
	} else if l.ch == '[' {
		tok = newToken(token.LBRACK, l.ch)
	} else if l.ch == ']' {
		tok = newToken(token.RBRACK, l.ch)
	} else if l.ch == ',' {
		tok = newToken(token.COMMA, l.ch)
	} else if l.ch == '"' {
		tok = newToken(token.QUOTES, l.ch)
	} else if l.ch == ':' {
		tok = newToken(token.DD, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{tokenType, string(ch)}
}
