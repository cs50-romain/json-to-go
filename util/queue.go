package queue

import "github.com/cs50-romain/jsontogo/internal"
import "fmt"

type Node struct {
	value token.Token
	next  *Node
}

type Queue struct {
	head *Node
}

func (q *Queue) Print() {
	if q.head == nil {
		return
	}

	curr := q.head
	for curr.next != nil {
		fmt.Println("Current node:", curr)
		curr = curr.next
	}
}

func (q *Queue) Pop() *token.Token {
	if q.head == nil {
		return &token.Token{"", ""} 
	}

	result := q.head
	q.head = q.head.next
	return &result.value
}

func (q *Queue) Peek() *token.Token {
	if q.head == nil {
		return &token.Token{"", ""} 
	}

	curr := q.head
	for curr.next != nil {
		curr = curr.next
	}
	return &curr.value
}

func (q *Queue) Push(s token.Token) {
	node := &Node{s, nil}
	if q.head == nil {
		q.head = node
		return
	}

	curr := q.head
	for curr.next != nil {
		curr = curr.next
	}
	curr.next = node
}

func Init() *Queue {
	return &Queue{nil}
}
