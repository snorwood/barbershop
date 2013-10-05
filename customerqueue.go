package barbershop

type CustomerQueue struct {
	queue *GeneralQueue
}

func NewCustomerQueue(capacity int) *CustomerQueue {
	customerQueue := new(CustomerQueue)
	customerQueue.queue = NewGeneralQueue(capacity)

	return customerQueue
}

func (self *CustomerQueue) Enqueue(customer *Customer) error {
	return self.queue.Enqueue(customer)
}

func (self *CustomerQueue) Dequeue() (*Customer, error) {
	value, err := self.queue.Dequeue()
	customer := value.(*Customer)

	return customer, err
}

func (self *CustomerQueue) Peek() (*Customer, error) {
	value, err := self.queue.Peek()
	customer := value.(*Customer)

	return customer, err
}

func (self *CustomerQueue) Size() int {
	return self.queue.Size()
}

func (self *CustomerQueue) Full() bool {
	return self.queue.Full()
}
