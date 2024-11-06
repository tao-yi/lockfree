package lockfree

import (
	"fmt"
	"strings"
	"sync"
)

type node struct {
	val  int
	next *node
}

type Stack1WithLock struct {
	top  *node
	mu   sync.Mutex
	size int
}

func NewStack1WithLock() *Stack1WithLock {
	return &Stack1WithLock{}
}

func (s *Stack1WithLock) Push(val int) *node {
	newNode := &node{val: val}
	s.mu.Lock()
	defer s.mu.Unlock()
	oldTop := s.top
	newNode.next = oldTop
	s.top = newNode
	s.size++
	return newNode
}

func (s *Stack1WithLock) String() string {
	cur := (*node)(s.top)
	var sb strings.Builder
	for cur != nil {
		sb.WriteString(fmt.Sprint(cur.val))
		sb.WriteString("->")
		cur = cur.next
	}
	if sb.Len() != 0 {
		sb.WriteString("nil")
	}
	return sb.String()
}
