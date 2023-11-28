package queue

import "github.com/cs50-romain/jsontogo/internal"

type Node struct {
	value token.Token
	next  *Node
}

type Queue struct {
	head *Node
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
	} else {
		curr := q.head
		for curr != nil {
			curr = curr.next
		}
		curr = node
	}
}

func Init() *Queue {
	return &Queue{nil}
}
