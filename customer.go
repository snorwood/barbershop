package barbershop

// Customer that frequents the barbershop for a close shave
type Customer struct {
	sendRequestChan chan Request
}

// NewCustomer initializes an instance of a customer
func NewCustomer(sendRequestChan chan Request) *Customer {
	customer := new(Customer)
	customer.sendRequestChan = sendRequestChan
	return customer
}

// SendRequest makes a request to the customers manager i.e. barbershop
func (self *Customer) SendRequest(req Request) chan Response {
	answer := make(chan Response)
	req.SetAnswerChannel(answer)
	self.sendRequestChan <- req
	return answer
}
