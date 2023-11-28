package main

import (
	"fmt"
	//"io"
	"bufio"
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cs50-romain/jsontogo/util"
	"github.com/cs50-romain/jsontogo/internal"
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

var tokens []token.Token
var tokenqueue *queue.Queue

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

func (t *Treeast) traversal(node *Node, closeobj bool, str *strings.Builder) string {
	if node == nil {
		return "" 
	}
	
	if node != t.head {
		//fmt.Println("node:", node.key, node.level, prevnodelvl)
		str.WriteString(turnToGo(node, closeobj))
	}
	
	for i := 0; i < len(node.children); i++ {
		if i > 0 && i == len(node.children)-1 {
			_ = t.traversal(node.children[i], true, str)
		} else {
			_ = t.traversal(node.children[i], false, str)
		}
	}
	return str.String()
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
		return 3
	} else if ch == RIGHT_BR {
		return 3
	} else if ch == QUOTE {
		return 1
	} else {
		return 0 
	}
}

func readInput(input string) []rune{
	chars := []rune{}

	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Error:",err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Append characters to array for simpler pass through
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println("Line: ", line)
		for _, char := range line {
			chars = append(chars, char)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return chars
}

// token.Tokenizer. Create an array of tokens for now
func lexer(input string){
	// Has to start as false otherwise values and identifiers will be swapped
	isIdentifier := false 
	isValue := false

	var charBuffer string
	chars := readInput(input)

	// token.Tokenizing 
	for i := 0; i < len(chars); i++ {
		char := chars[i]
		useBuffer := readChar(char)

		if useBuffer == 3 /*Array*/ {
			// Everything inside the array is a value
			if string(char) == "[" {
				isValue = true
			} else {
				isValue = false
				isIdentifier = true
			}

			token := token.Token{
				Lexeme: "Array",
				Value: string(char),
			}
			tokens = append(tokens, token)
			tokenqueue.Push(token)
		} else if useBuffer == 2 /**/{
			if string(char) == "{" {
				if !isIdentifier {
					isIdentifier = true
				} else {
					isIdentifier = false
				}
			}

			token := token.Token{
				Lexeme: "startObject",
				Value: string(char),
			}
			tokens = append(tokens, token)
		} else if useBuffer == 1 {
			if useBuffer == 2 {
			}
			if len(charBuffer) > 0 {
				if charBuffer[0] == ' ' {
					charBuffer = ""
				} else {
					token := token.Token{
						Lexeme: "",
						Value: charBuffer,
					}
					if isValue {
						token.Lexeme = "value"
					} else if isIdentifier {
						token.Lexeme = "identifier"
						isIdentifier = false
					} else {
						token.Lexeme = "value"
						isIdentifier = true
					}
					if !isWhiteSpace(charBuffer ){
						tokens = append(tokens, token)
						tokenqueue.Push(token)
					}
					charBuffer = ""
				}
			}

			token := token.Token{
				Lexeme: "punctuator",
				Value: string(char),
			}
			tokens = append(tokens, token)
			tokenqueue.Push(token)
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

	for i := 0; i < len(tokens)-1; i++ {
		token := tokens[i]
		fmt.Println(token)

		if token.Value == "{" {
			level++
		} else if token.Value == "}" {
			level--
		}


		if token.Lexeme == "Array" {
			
		} else if token.Lexeme == "value" {
			if node.value != "" || node.value == "array" {
				node.value = "array"
			} else {
				node.value = token.Value
			}
		} else if token.Lexeme == "startObject" {
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
		} else if token.Lexeme == "identifier" {
			node = &Node{token.Value, "", level, nil}
			curr.children = append(curr.children, node)
		}
	}
}

func strToByte(input []string) []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(input)
	return buf.Bytes()
}

// Generate a go struct from the AST tree
// Use NLR traversal
func toGo(tree *Treeast) string {
	var b strings.Builder
	b.WriteString("type Object struct{\n")
	_ = tree.traversal(tree.head, false, &b)
	
	return b.String()
}

func turnToGo(node *Node, closeobj bool) string {
	tab := strings.Repeat("\t", node.level)
	result := ""

	if node.value == "struct" {
		result = tab + node.key + "\t" + "struct {\n"
	} else if node.value == "array" {		
		if len(node.value) < 20  {
			result = tab + node.key + "\t\t" + "[]string\n"
		} else {
			result = tab + node.key + "\t" + "[]string\n"
		}
	}else {
		if len(node.value) < 20  {
			result = tab + node.key + "\t\t" + "string\n"
		} else {
			result = tab + node.key + "\t" + "string\n"
		}
	}

	if closeobj == true {
		tab = strings.Repeat("\t", node.level-1)
		result += tab + "}\n"
	}

	return result
}

func toClipboard(output string) {
	var copyCmd *exec.Cmd

	copyCmd = exec.Command("xclip", "-selection", "C")

	in, err := copyCmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	if err := copyCmd.Start(); err != nil {
		log.Fatal(err)
	}

	if _, err := in.Write([]byte(output)); err != nil {
		log.Fatal(err)
	}

	if err := in.Close(); err != nil {
		log.Fatal(err)
	}

	copyCmd.Wait()
}

func parseCmd(flags []string) {
	doCopy := false
	var output string 
	if len(flags) > 0 {
		for idx, flag := range flags {
			if flag == "jsontogo.go" || idx == 0 {
				continue
			} else if flag == "copy" {
				doCopy = true
			} else if len(flag) > 5 && flag[len(flag)-5:len(flag)] == ".json" {
				tree := InitAST()
				root := &Node{"root", "root", 1, nil}
				tree.head = root
				lexer(flag)
				parser(tree)
				output = toGo(tree)
			} else {
				fmt.Println("Invalid command", flag)
			}
		}
	}

	if doCopy {
		toClipboard(output)
		fmt.Println("[+] Output copied to clipboard")
	} else {
		fmt.Println(output)
	}
}

func main() {
	tokenqueue = queue.Init()
	parseCmd(os.Args)
}
