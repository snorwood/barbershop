package barbershop

import (
// "mathrand"
)

// Customer that frequents the barbershop for a close shave
type Customer struct {
	sendRequestChan    chan Request
	receiveRequestChan chan Request
}

// NewCustomer initializes an instance of a customer
func NewCustomer(sendRequestChan chan Request) *Customer {
	customer := new(Customer)
	customer.sendRequestChan = sendRequestChan
	customer.receiveRequestChan = make(chan Request)
	return customer
}

// SendRequest makes a request to the customers manager i.e. barbershop
func (self *Customer) SendRequest(req Request) chan Subscriber {
	answer := make(chan Subscriber)
	req.SetAnswerChannel(answer)
	self.sendRequestChan <- req
	return answer
}

func (self *Customer) GetReceiveRequestChan() chan Request {
	return self.receiveRequestChan
}

func (self *Customer) Start() {
	request := NewRequest("Barber", "HELLO")
	barber := <-self.SendRequest(request)

	select {
	case message := <-barber.Send:
		print(message)
	case <-barber.StopReceiving:
		print("Done")
	}

}
