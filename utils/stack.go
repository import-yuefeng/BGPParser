// MIT License

// Copyright (c) 2019 Yuefeng Zhu

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
