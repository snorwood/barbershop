package barbershop

type Customer struct {
	sendRequestChan chan Request
	recieveAnswer   chan string
}

func NewCustomer(sendRequestChan chan Request) *Customer {
	customer := new(Customer)
	customer.sendRequestChan = sendRequestChan

	customer.recieveAnswer = make(chan string)

	return customer
}

func (self *Customer) SendRequest(req Request) chan string {
	answer := make(chan string)

	req.SetAnswerChannel(answer)
	self.sendRequestChan <- req

	return answer
}
