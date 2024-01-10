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
	} else if l.ch == ':' && !isLetter(l.input[l.position - 1]){
		tok = newToken(token.DD, l.ch)
	} else {
		// Check if its a string or a digit or a key
		if isLetter(l.ch) {
			tok.Literal = l.readKey()
			tok.Type = token.KEY
			return tok
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readKey() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '/' || ch == '\\'
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
	return token.Token{Type: tokenType, Literal: string(ch)}
}
