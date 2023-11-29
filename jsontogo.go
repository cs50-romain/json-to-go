package main

import (
	"fmt"
	//"io"
	"bufio"
	//"bytes"
	//"encoding/gob"
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
	ttype		string
	key		string
	value		[]string
	level		int
	children	[]*Node
}

type Treeast struct {
	head	*Node
}

func InitAST() *Treeast{
	return &Treeast{nil}
}

func (t *Treeast) traversal(node *Node, closeobj bool, str *strings.Builder) string {
	if node == nil {
		return "" 
	}
	
	fmt.Printf("node: %s %s %s level: %d\n", node.ttype, node.key, node.value, node.level)
	str.WriteString(turnToGo(node, closeobj))
	
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
				Lexeme: "object",
				Value: string(char),
			}
			tokens = append(tokens, token)
			tokenqueue.Push(token)
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

func parser(tree *Treeast, node *Node, prevIdentifierNode *Node, level int) {
	if tree.head == nil {
		root := tokenqueue.Pop()
		if root.Lexeme == "object" && root.Value == "{" {
			tree.head = &Node{"Object", "", []string{root.Value}, level, []*Node{}}	
		} else {
			tree.head = &Node{"Array", "", []string{root.Value}, level, []*Node{}}
		}
		node = tree.head
		level++
	}
	fmt.Println(level)

	// BUILD TEST CASE FOR RECURSION
	if tokenqueue.Peek() == nil {
		fmt.Println("Returning: next token is nil", tokenqueue.Peek())
		return
	}

	if tokenqueue.Peek().Value == "}" {
		_ = tokenqueue.Pop()
		return
	}

	if tokenqueue.Peek().Value == "]" {
		fmt.Println("Returning: next token is a RB", tokenqueue.Peek())
		_ = tokenqueue.Pop()
		return
	}

	currnode := node
	
	for tokenqueue.Peek() != nil {
		nextoken := tokenqueue.Pop()


		if nextoken == nil {
			return
		}

		if nextoken.Value == "}" || nextoken.Value == "]" {
			level--
		} else if nextoken.Value == "{" || nextoken.Value == "[" {
			level++
		}

		// Append to current node's children and recurse
		if nextoken.Lexeme == "object" { // The token is
			currnode.children = append(currnode.children, &Node{"Object", "", []string{nextoken.Value}, level, []*Node{}})

			parser(tree, currnode.children[len(currnode.children) - 1], prevIdentifierNode, level)
		} else if nextoken.Lexeme == "identifier" {
			newnode := &Node{"Property", nextoken.Value, []string{}, level, []*Node{}}
			currnode.children = append(currnode.children, newnode)
			// PEEK and POP until we peek either a value or an object or an array
			prevIdentifierNode = newnode
		} else if nextoken.Lexeme == "Array" {
			currnode.children = append(currnode.children, &Node{"Array", "", []string{nextoken.Value}, level, []*Node{}})
		} else if nextoken.Lexeme == "value" {
			if prevIdentifierNode != nil {
				prevIdentifierNode.value = append(prevIdentifierNode.value, nextoken.Value)
			}
		}
	}
	return
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

	if node.value != nil/*== "struct"*/ {
		result = tab + node.key + "\t" + "struct {\n"
	} else if node.value != nil /* == "array" */{		
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
		tab = strings.Repeat("\t", node.level)
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
				lexer(flag)
				parser(tree, tree.head, nil, 0)
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
		fmt.Println("Output:", output)
	}
}

func main() {
	tokenqueue = queue.Init()
	parseCmd(os.Args)
}
