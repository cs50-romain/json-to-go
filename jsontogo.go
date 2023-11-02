package main

import (
	"fmt"
	//"io"
	"bufio"
	"os"
)

// Possible Characters
var (
	LEFT_CB = '{'  //0
	RIGHT_CB = '}' //1
	LEFT_PA = '('  //2
	RIGHT_PA = ')' //3
	LEFT_BR = '['  //4
	RIGHT_BR = ']' //5
	DD = ':'       //6
	COMA = ','     //7
	QUOTE = '"'    //8
	//CHAR = '[a-zA-Z]'
	W_SPACE = ' '
)

//Tokens:
// Identifier start with " and ends with " and is before the :
// String starts with " and ends with " but is always after a :
// Punctuator is either { } : [ ] ,

type Token struct {
	lexeme string //String
	value string
}

var tokens []Token

// Tokenizer. Create an array of tokens for now
func lexer() {
	chars := []rune{}
	var charBuffer string

	file, err := os.Open("sample.json")
	if err != nil {
		fmt.Println("Error:",err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Append characters to array for simpler pass through
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Line: ", line)
		for _, char := range line {
			chars = append(chars, char)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Tokenizing 
	for i := 0; i < len(chars); i++ {
		char := chars[i]
		useBuffer := readChar(char)

		if !useBuffer {
			if len(charBuffer) > 0 {
				token := Token{
					lexeme: "string",
					value: charBuffer,
				}
				if !isWhiteSpace(charBuffer) {
					tokens = append(tokens, token)
				}
				charBuffer = ""
			}

			token := Token{
				lexeme: "punctuator",
				value: string(char),
			}
			tokens = append(tokens, token)
		} else {
			charBuffer += string(char)			
		}
	}

	fmt.Println(tokens)
}

func isWhiteSpace(str string) bool {
	if str == " " {
		return true
	}
	return false
}

//read character
func readChar(ch rune) bool {
	if ch == W_SPACE {
		fmt.Println(" whitespace")
		return true 
	} else if ch == LEFT_CB {
		fmt.Println(" left curly bracket")
		return false
	} else if ch == RIGHT_CB {
		fmt.Println(" right curly bracket")
		return false
	} else if ch == LEFT_PA {
		fmt.Println(" left parentheses")
		return false
	} else if ch == RIGHT_PA {
		fmt.Println(" right parentheses")
		return false
	} else if ch == DD {
		fmt.Println(" double dot")
		return false
	} else if ch == COMA {
		fmt.Println(" comma")
		return false
	} else if ch == LEFT_BR {
		fmt.Println(" left bracket")
		return false
	} else if ch == RIGHT_BR {
		fmt.Println(" right bracket")
		return false
	} else if ch == QUOTE {
		fmt.Println(" quote")
		return false
	} else {
		fmt.Println(" character")
		return true 
	}
}

func parser() {

}

func main() {
	lexer()	
}
