package barbershop

type Request struct {
	answer  chan string
	message string
}

func NewRequest(message string) Request {
	request := Request{}

	request.message = message

	request.answer = make(chan string)

	return request
}

func (self *Request) GetAnswerChannel() chan string {
	return self.answer
}

func (self *Request) GetMessage() string {
	return self.message
}

func (self *Request) SetAnswerChannel(answer chan string) {
	self.answer = answer
}
