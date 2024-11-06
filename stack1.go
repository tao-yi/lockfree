package lockfree

import (
	"sync/atomic"
	"unsafe"
)

type Stack1 struct {
	top  unsafe.Pointer
	size atomic.Int32
}

type element struct {
	value int
	next  unsafe.Pointer
}

func NewStack1() *Stack1 {
	return &Stack1{}
}

func (s *Stack1) Push(val int) unsafe.Pointer {
	newNode := &element{value: val}
	for {
		oldTop := atomic.LoadPointer(&s.top)
		newNode.next = oldTop
		if atomic.CompareAndSwapPointer(&s.top, oldTop, unsafe.Pointer(newNode)) {
			s.size.Add(1)
			return unsafe.Pointer(newNode)
		}
	}
}

func (s *Stack1) Pop() (int, bool) {
	for {
		oldTop := atomic.LoadPointer(&s.top)
		if oldTop == nil {
			return 0, false
		}

		newTop := (*element)(oldTop).next
		if atomic.CompareAndSwapPointer(&s.top, oldTop, unsafe.Pointer(newTop)) {
			s.size.Add(-1)
			return (*element)(oldTop).value, true
		}
	}
}

func (s *Stack1) Len() int32 {
	return s.size.Load()
}
