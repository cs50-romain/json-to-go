package main

import (
	"fmt"
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

// Add checking for integer (currently lexer will only recognize strings "something" as a value to an identifier

// Possible Characters
const (
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
var nodesarray []*Node

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

// Maybe add a function that checks if a another node is equal (same key, level) to this node. If so 
// To make it faster, go down the tree by the level of the node(a child is 1 level)

func (t *Treeast) traversal(node *Node, closeobj bool, str *strings.Builder) string {
	if node == nil {
		return "" 
	}
	
	//fmt.Printf("node{type: %s | key: %s | value: %s | level: %d}\n", node.ttype, node.key, node.value, node.level)
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
	buf := []byte{}

	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Error:",err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 2048*1024)

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

func lexer(input string){
	// Has to start as false otherwise values and identifiers will be swapped
	isIdentifier := false 
	isValue := false

	var charBuffer string
	var prevChar rune
	chars := readInput(input)

	// token.Tokenizing 
	for i := 0; i < len(chars); i++ {
		char := chars[i]
		useBuffer := readChar(char)

		if useBuffer == 3 /*Array*/ {
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
		} else if useBuffer == 2 /*Curly B*/{
			if string(char) == "{" {
				isIdentifier = true
				isValue = false
			} else {
				isIdentifier = false
			}

			token := token.Token{
				Lexeme: "object",
				Value: string(char),
			}
			tokens = append(tokens, token)
			tokenqueue.Push(token)
		} else if useBuffer == 1 {
			if prevChar != '"' && char == ':' {
				charBuffer += string(char)
			} else {
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
			}
		} else {
			charBuffer += string(char)			
		}
		prevChar = char
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

	if tokenqueue.Peek() == nil {
		fmt.Println("Returning: next token is nil", tokenqueue.Peek())
		return
	}

	if tokenqueue.Peek().Value == "}" {
		fmt.Println("Returning: next token is CB", tokenqueue.Peek())
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

		if nextoken.Lexeme == "object" {
			newnode := &Node{"Object", "", []string{nextoken.Value}, level, []*Node{}}
			if !doesExist(newnode) {
				nodesarray = append(nodesarray, newnode)
				currnode.children = append(currnode.children, newnode)
				parser(tree, currnode.children[len(currnode.children) - 1], prevIdentifierNode, level)
			}
		} else if nextoken.Lexeme == "identifier" {
			newnode := &Node{"Property", nextoken.Value, []string{}, level, []*Node{}}
			if !doesExist(newnode) {
				nodesarray = append(nodesarray, newnode)
				currnode.children = append(currnode.children, newnode)
				if tokenqueue.Peek().Value == "\"" {
					_ = tokenqueue.Pop()
					if tokenqueue.Peek().Value == ":" {
						_ = tokenqueue.Pop()
						// Add dummy values for now
						if tokenqueue.Peek().Lexeme == "Array" {
							newnode.value = append(newnode.value, "a", "b")
						} else {
							newnode.value = append(newnode.value, "a")
						}
					}
				}
				//prevIdentifierNode = newnode
			}
		} else if nextoken.Lexeme == "Array" {
			newnode := &Node{"Array", "", []string{nextoken.Value}, level, []*Node{}}
			if !doesExist(newnode) {
				nodesarray = append(nodesarray, newnode)
				currnode.children = append(currnode.children, newnode)
			}
		} else if nextoken.Lexeme == "value" {
			/*
			fmt.Println("Previous node:", prevIdentifierNode)
			fmt.Print("\n\n\n")
			if prevIdentifierNode != nil {
				prevIdentifierNode.value = append(prevIdentifierNode.value, nextoken.Value)
			}
			*/
		}
	}
	return
}

// Generate a go struct from the AST tree
// Use NLR traversal
func toGo(tree *Treeast) string {
	var b strings.Builder
	b.WriteString("type AutomatedType struct {\n")
	_ = tree.traversal(tree.head, false, &b)
	b.WriteString("}")
	
	return b.String()
}

func turnToGo(node *Node, closeobj bool) string {
	tab := strings.Repeat("\t", node.level)
	result := ""

	if node.ttype == "Object" /*== "struct"*/ {
	} else if node.ttype == "Array" /* == "array" */{		
	} else if node.ttype == "Property" && len(node.value) >= 1 {
		if len(node.value) > 1 {
			if len(node.key) < 10  {
				result = tab + node.key + "\t\t" + "[]string\n"
			} else {
				result = tab + node.key + "\t" + "[]string\n"
			}
		} else {
			if len(node.key) < 10 {
				result = tab + node.key + "\t\t" + "string\n"
			} else {
				result = tab + node.key + "\t" + "string\n"
			}
		}
	} else {
		result = tab + node.key + "\t" + "struct {\n"
	}

	if closeobj == true {
	}

	return result
}

func doesExist(needlenode *Node) bool {
	for _, node := range nodesarray {
		if needlenode.level == node.level && needlenode.key == node.key && needlenode.ttype == node.ttype {
			return true
		}
	}
	return false
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

/*
Fixing multiples in one array object:
Solution 1:
Create a map, if exact node already exist (check to ttype, key, level) then check if value is same. If different, add the value of current node to already existing node.

It kinda works but not good enough. Bigger data resutls in errors. Fixed solution isn't the way.

Solution 2:
Create a single token Array that holds anything in between [ ]. Except if the object starts with an array.
*/
