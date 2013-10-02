package barbershop

import (
	"fmt"
)

type node struct {
	prev  *interface{}
	value interface{}
	next  *interface{}
}

type GeneralQueue struct {
	head     node
	tail     node
	size     int
	capacity int
}

func (self *GeneralQueue) Enqueue(value interface{}) error {
	if self.size < self.capacity {
		newNode := Node{
			next:  null,
			prev:  null,
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
			head = newNode
		}
		size++

		return nil
	}

	return QueueOverflowException{index: size, size: size}
}

func (self *GeneralQueue) Dequeue() (interface{}, err) {
	if self.size > 0 {
		frontNode := self.head.next
		head.next = frontNode.next
		size--

		return frontNode.value
	}
}

type QueueOverflowException struct {
	index int
	size  int
	err   string
}

func (self QueueOverflowException) Error() string {
	return fmt.Sprintf("QueueOverflowException: Index (%d) out of range (%d).", self.index, self.size)
}
