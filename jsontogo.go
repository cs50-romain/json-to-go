package main

import (
	"fmt"
	//"io"
	"bufio"
	"os"
)

// Possible Characters
var (
	LEFT_BR = '{'
	RIGHT_BR = '{'
	LEFT_PA = '('
	RIGHT_PA = ')'
	DD = ':'
	COMA = ','
	QUOTE = '"'
	//CHAR = '[a-zA-Z]'
	W_SPACE = ' '
)

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
	if ch == LEFT_BR {
		fmt.Println(" left bracket")
	} else if ch == RIGHT_BR {
		fmt.Println(" left bracket")
	} else if ch == LEFT_PA {
		fmt.Println(" left parentheses")
	} else if ch == RIGHT_PA {
		fmt.Println(" right parentheses")
	} else if ch == DD {
		fmt.Println(" double dot")
	} else if ch == COMA {
		fmt.Println(" comma")
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
