package main

import (
	"fmt"
	//"io"
	"bufio"
	"os"
	"strings"
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

var isIdentifier bool
var tokens []Token

//Tokens:
// Identifier start with " and ends with " and is before the :
// String starts with " and ends with " but is always after a :
// Punctuator is either { } : [ ] ,

type Token struct {
	lexeme	string //String
	value	string
}

// AST Node
type Node struct {
	key	string
	value	string
	level	int
	children []*Node
}

type Treeast struct {
	head	*Node
}

func InitAST() *Treeast{
	return &Treeast{nil}
}

func (t *Treeast) print() {
	if t.head == nil {
		return
	}

	curr := t.head
	for _, child := range curr.children {
		fmt.Println(child.key, child.value, child.level)
	}
}

func (t *Treeast) traversal(node *Node, closeobj bool) {
	if node == nil {
		return
	}
	
	if node != t.head {
		//fmt.Println("node:", node.key, node.level, prevnodelvl)
		turnToGo(node, closeobj)
	}
	
	for i := 0; i < len(node.children); i++ {
		if i > 0 && i == len(node.children)-1 {
			t.traversal(node.children[i], true)
		} else {
			t.traversal(node.children[i], false)
		}
	}
}

func isWhiteSpace(str string) bool {
	if str == " " {
		return true
	}
	return false
}

//read character
func readChar(ch rune) int {
	if ch == W_SPACE {
		return 0 
	} else if ch == LEFT_CB {
		return 2 
	} else if ch == RIGHT_CB {
		return 2
	} else if ch == LEFT_PA {
		return 1
	} else if ch == RIGHT_PA {
		return 1
	} else if ch == DD {
		return 1
	} else if ch == COMA {
		return 1
	} else if ch == LEFT_BR {
		return 1
	} else if ch == RIGHT_BR {
		return 1
	} else if ch == QUOTE {
		return 1
	} else {
		return 0 
	}
}

// Tokenizer. Create an array of tokens for now
func lexer(input string){

	// Has to start as false otherwise values and identifiers will swapped
	isIdentifier = false 
	chars := []rune{}
	var charBuffer string

	file, err := os.Open(input)
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

		if useBuffer == 2 {
			if string(char) == "{" {
				if !isIdentifier {
					isIdentifier = true
				} else {
					isIdentifier = false
				}
			}

			token := Token{
				lexeme: "startObject",
				value: string(char),
			}
			tokens = append(tokens, token)
		} else if useBuffer == 1 {
			if useBuffer == 2 {
				fmt.Println("Embedding")
			}
			if len(charBuffer) > 0 {
				if charBuffer[0] == ' ' {
					charBuffer = ""
				} else {
					token := Token{
						lexeme: "",
						value: charBuffer,
					}
					if isIdentifier {
						token.lexeme = "identifier"
						isIdentifier = false
					} else {
						token.lexeme = "value"
						isIdentifier = true
					}
					if !isWhiteSpace(charBuffer ){
						tokens = append(tokens, token)
					}
					charBuffer = ""
				}
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
}

// Going through token array and create an AST tree by creating nodes
func parser(tree *Treeast) {
	level := 1 
	// Start helps to distinguish whether it is the start of a json object; true = start of json object
	start := true
	var node *Node
	curr := tree.head
	var prevnodes []*Node

	for i := 1; i < len(tokens)-1; i++ {
		token := tokens[i]

		if token.value == "{" {
			level++
		} else if token.value == "}" {
			level--
		}

		if token.lexeme == "value" {
			node.value = token.value
		} else if token.lexeme == "startObject" {
			if start == true {	
				node.value = "struct"
				prevnodes = append(prevnodes, curr)
				curr = node
				start = false
			} else {
				curr = prevnodes[0]
				if len(prevnodes) > 0 {
					prevnodes = append(prevnodes[1:])
				} else {
					prevnodes = []*Node{}
				}
				start = true
			}
		} else if token.lexeme == "identifier" {
			node = &Node{token.value, "", level, nil}
			curr.children = append(curr.children, node)
		}
	}
}

// Generate a go struct from the AST tree
// Use NLR traversal
func toGo(tree *Treeast) {
	fmt.Println("type Object struct {")
	tree.traversal(tree.head, false)
}

func turnToGo(node *Node, closeobj bool) {
	tab := strings.Repeat("\t", node.level)

	if node.value == "struct" {
		fmt.Printf("%s%s\t%s\n", tab, node.key, "struct {")
	} else {
		if len(node.value) < 20  {
			fmt.Printf("%s%s\t\t%s\n", tab, node.key, "string")
		} else {
			fmt.Printf("%s%s\t%s\n", tab, node.key, "string")
		}
	}

	if closeobj == true {
		tab = strings.Repeat("\t", node.level-1)
		fmt.Printf("%s}\n", tab)
	}
}

func parseCmd(flags []string) {
	if flags[1] != "copy" {
		file := flags[1]
		tree := InitAST()
		root := &Node{"root", "root", 1, nil}
		tree.head = root
		lexer(file)
		parser(tree)
		toGo(tree)

	}
	if flags[2] == "copy" {
		fmt.Println("Object copied to clipboard")
	}
}

func main() {
	parseCmd(os.Args)
}

// TODOD:
// Create a queue to hold tokens; parser can pop and peek
