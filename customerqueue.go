package barbershop

type chair struct {
	Occupant Customer
	Next     *Customer
	Prev     *Customer
}

type CustomerQueue struct {
	queue GeneralQueue
}

// func NewLineUp(capacity int) {
// 	capacity = capacity
// 	head = new(chair)
// 	head.Next = null
// 	head.prev = null
// 	size = 0
// }

func (self *CustomerQueue) Enqueue(customer Customer) error {
	return self.queue.Enqueue(customer)
}

func (self *CustomerQueue) Dequeue() (Customer, error) {
	value, err := self.queue.Dequeue()
	customer := value.(Customer)

	return customer, err
}

func (self *CustomerQueue) Peek() (Customer, error) {
	value, err := self.peek.Dequeue()
	customer := value.(Customer)

	return customer, err
}
