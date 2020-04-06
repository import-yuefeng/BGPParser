package utils

import "errors"

type Queue struct {
	val         []interface{}
	front, rare int
	size        int
	capacity    int
}

func NewQueue() *Queue {
	// default size: 2
	return &Queue{
		val:      make([]interface{}, 2),
		front:    0,
		rare:     0,
		size:     0,
		capacity: 2,
	}
}

func (q *Queue) Push(x interface{}) {
	if q.size == q.capacity {
		buf := NewQueue()
		buf.val = make([]interface{}, q.capacity*2)
		size := q.size
		for i := 0; i < size; i++ {
			buf.val[i] = q.val[q.front]
			q.front = (q.front + 1) % q.capacity
		}
		q.val = buf.val
		q.front = 0
		q.rare = q.size
		q.capacity *= 2
	}
	q.val[q.rare] = x
	q.rare = (q.rare + 1) % q.capacity
	q.size++
	return
}

func (q *Queue) Pop() (x interface{}, err error) {
	if q.IsEmpty() {
		return nil, errors.New("Queue is empty")
	}
	x = q.val[q.front]
	q.front = (q.front + 1) % q.capacity
	q.size--
	return x, nil
}

func (q *Queue) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue) Size() int {
	return q.size
}
