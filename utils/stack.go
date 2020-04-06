package utils

import (
	"errors"
)

type Stack struct {
	val      []interface{}
	top      int
	size     int
	capacity int
}

func NewStack() *Stack {
	// default size: 2
	return &Stack{
		val:      make([]interface{}, 2),
		top:      -1,
		size:     0,
		capacity: 2,
	}
}

func (s *Stack) Push(x interface{}) {
	if s.size == s.capacity {
		buf := NewStack()
		buf.val = make([]interface{}, s.capacity*2)
		copy(buf.val, s.val)
		s.val = buf.val
		s.capacity *= 2
	}
	s.top++
	s.val[s.top] = x
	s.size++
	return
}

func (s *Stack) Pop() (x interface{}, err error) {
	if s.size == 0 {
		return nil, errors.New("Stack is empty!")
	}
	x = s.val[s.top]
	s.top--
	s.size--
	return x, nil
}

func (s *Stack) Reset() {
	s.val = s.val[:0]
	s.top, s.size, s.capacity = -1, 0, 0
	return
}

func (s *Stack) IsEmpty() bool {
	return s.size == 0
}

func (s *Stack) Top() (x interface{}) {
	if s.size == 0 {
		return nil
	}
	return s.val[s.top]
}

func (s *Stack) Size() int {
	return s.size
}
