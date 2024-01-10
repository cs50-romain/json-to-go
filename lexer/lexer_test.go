package lexer

import (
	"fmt"
	"testing"

	"github.com/cs50-romain/jsontogo/token"
)

func TestNextToken(t *testing.T) {
	input := `{}[],":`

	tests := []struct{
		expectedType	token.TokenType
		expectedLiteral	string
	}{
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACK, "["},
		{token.RBRACK, "]"},
		{token.COMMA, ","},
		{token.QUOTES, string('"')},
		{token.DD, ":"},
		{token.KEY, "Hello"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		fmt.Println(i, tt, tok)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
