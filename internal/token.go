package token

//Tokens:
// Identifier start with " and ends with " and is before the :
// String starts with " and ends with " but is always after a :
// Punctuator is either { } : [ ] ,

type Token struct {
	Lexeme	string //String
	Value	string
}
