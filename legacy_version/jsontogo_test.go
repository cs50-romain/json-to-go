package main

import "testing"
import "strings"

func BenchmarkTraversal(b *testing.B) {
	node1 := &Node{"string", "value2", 2 /*level*/, nil}
	node2 := &Node{"string", "value3", 2 /*level*/, nil}
	node3 := &Node{"string", "value4", 2 /*level*/, nil}
	head := &Node{"struct", "value1", 1 /*level*/, nil}
	head.children = append(head.children, node1, node2, node3)
	var bs strings.Builder

	tree := InitAST()
	tree.head = head

	_ = tree.traversal(head, false, bs)
}
