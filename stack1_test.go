package lockfree

import (
	"math/rand/v2"
	"sync"
	"testing"
)

func BenchmarkCASPush(b *testing.B) {
	s := NewStack1()
	for i := 0; i < b.N; i++ {
		s.Push(101)
	}
}

func BenchmarkLockPush(b *testing.B) {
	s := NewStack1WithLock()
	for i := 0; i < b.N; i++ {
		s.Push(101)
	}
}

func BenchmarkStack2Push(b *testing.B) {
	s := NewStack2[int]()
	for i := 0; i < b.N; i++ {
		s.Push(101)
	}
}

func BenchmarkParallelCASPush(b *testing.B) {
	s := NewStack1()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Push(101)
		}
	})
}

func BenchmarkParallelLockPush(b *testing.B) {
	s := NewStack1WithLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Push(101)
		}
	})
}

func BenchmarkParallelStack2(b *testing.B) {
	s := NewStack2[int]()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			s.Push(101)
		}
	})
}

func TestStack1CASPush(t *testing.T) {
	var wg sync.WaitGroup
	s := NewStack1()
	size := rand.IntN(1000)
	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Push(i)
		}()
	}

	wg.Wait()
	if int(s.size.Load()) != size {
		t.Errorf("size is wrong got %d, expected %d", s.size.Load(), size)
	}
}
func TestStack2CASPush(t *testing.T) {
	var wg sync.WaitGroup
	s := NewStack2[int]()
	size := rand.IntN(1000)
	for i := range size {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Push(i)
		}()
	}

	wg.Wait()
	if int(s.size.Load()) != size {
		t.Errorf("size is wrong got %d, expected %d", s.size.Load(), size)
	}
}
