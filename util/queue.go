package queue

type Node struct {
	value string
	next  *Node
}

type Queue struct {
	head *Node
}

func (q *Queue) Pop() string {
	if q.head == nil {
		return ""
	}

	result := q.head
	q.head = q.head.next
	return result
}

func (q *Queue) Peek() string {
	if q.head == nil {
		return ""
	}

	curr := q.head
	for curr.next != nil {
		curr = curr.next
	}
	return curr
}

func (q *Queue) Push(s string) {
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
