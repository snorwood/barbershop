package barbershop

type chair struct {
	Occupant Customer
	Next     *Customer
	Prev     *Customer
}

type LineUp struct {
	queue GeneralQueue
}

// func NewLineUp(capacity int) {
// 	capacity = capacity
// 	head = new(chair)
// 	head.Next = null
// 	head.prev = null
// 	size = 0
// }

func (self *LineUp) Enqueue(customer Customer) {

}
