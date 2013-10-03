package barbershop

import (
	"fmt"
)

type node struct {
	prev  *node
	value interface{}
	next  *node
}

type GeneralQueue struct {
	head     *node
	tail     *node
	size     int
	capacity int
}

func NewGeneralQueue(capacity int) *GeneralQueue {
	generalQueue := GeneralQueue{
		head: nil,
		tail: nil,
		size: 0,
		capacity: capacity
	}

	return &GeneralQueue
}

func (self *GeneralQueue) Enqueue(value interface{}) error {
	if self.size < self.capacity {
		newNode := &node{
			next:  nil,
			prev:  nil,
			value: value,
		}

		visitor := self.head

		for index := 0; index < self.size; index++ {
			visitor = visitor.next
		}

		visitor.next = newNode
		newNode.prev = visitor
		self.tail = newNode

		if self.size == 0 {
			self.head = newNode
		}
		self.size++

		return nil
	}

	return QueueOverflow{index: self.size, size: self.size}
}

func (self *GeneralQueue) Dequeue() (interface{}, error) {
	if self.size > 0 {
		frontNode := self.head.next
		self.head.next = frontNode.next
		self.size--

		if self.size == 0 {
			self.tail = nil
		}

		return frontNode.value, nil
	}

	return nil, EmptyQueue{}
}

func (self *GeneralQueue) Peek() (interface{}, error) {
	if self.size > 0 {
		return self.head.next.value, nil
	}

	return nil, EmptyQueue{}
}

type QueueOverflow struct {
	index int
	size  int
}

func (self QueueOverflow) Error() string {
	return fmt.Sprintf("QueueOverflow: Index (%d) out of range (%d).", self.index, self.size)
}

type EmptyQueue struct{}

func (self EmptyQueue) Error() string {
	return fmt.Sprintf("EmptyQueue: Dequeue attempted from already empty queue.")
}
