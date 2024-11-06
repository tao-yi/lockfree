package lockfree

import (
	"sync/atomic"
	"unsafe"
)

type Stack2[T any] struct {
	top  atomic.Pointer[element2[T]]
	size atomic.Int32
}

type element2[T any] struct {
	value T
	next  atomic.Pointer[element2[T]]
}

func NewStack2[T any]() *Stack2[T] {
	return &Stack2[T]{
		size: atomic.Int32{},
	}
}

func (s *Stack2[T]) Push(val T) unsafe.Pointer {
	newNode := &element2[T]{value: val}
	for {
		oldTop := s.top.Load()
		newNode.next.Store(oldTop)
		if s.top.CompareAndSwap(oldTop, newNode) {
			s.size.Add(1)
			return unsafe.Pointer(newNode)
		}
	}
}

func (s *Stack2[T]) Pop() (T, bool) {
	for {
		oldTop := s.top.Load()
		if oldTop == nil {
			var res T
			return res, false
		}

		newTop := oldTop.next.Load()
		if s.top.CompareAndSwap(oldTop, newTop) {
			s.size.Add(-1)
			return newTop.value, true
		}
	}
}

func (s *Stack2[T]) Len() int32 {
	return s.size.Load()
}
