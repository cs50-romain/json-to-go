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

type String struct {
	type string //String
	value string
}

type Punctuator struct {
	type string // Punctuator
	value rune
}

// Tokenizer. Create an array of tokens for now
func lexer() {
	reader()
}

// Read from file
func reader() {
	file, err := os.Open("sample.json")
	if err != nil {
		fmt.Println("Error:",err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Line: ", line)
		readLine(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func readLine(line string){
	for _, char := range line {
		fmt.Print("Char:", string(char))
		readChar(char)
	}
}

//read character
func readChar(ch rune) {
	if ch == W_SPACE {
		fmt.Println(" whitespace")
	} else if ch == LEFT_CB {
		fmt.Println(" left curly bracket")
	} else if ch == RIGHT_CB {
		fmt.Println(" right curly bracket")
	} else if ch == LEFT_PA {
		fmt.Println(" left parentheses")
	} else if ch == RIGHT_PA {
		fmt.Println(" right parentheses")
	} else if ch == DD {
		fmt.Println(" double dot")
	} else if ch == COMA {
		fmt.Println(" comma")
	} else if ch == LEFT_BR {
		fmt.Println(" left bracket")
	} else if ch == RIGHT_BR {
		fmt.Println(" right bracket")
	} else if ch == QUOTE {
		fmt.Println(" quote")
	} else {
		fmt.Println(" character")
	}
}

func parser() {

}

func main() {
	lexer()	
}
