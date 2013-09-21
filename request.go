package barbershop

// Request is used to make requests from participants to the manager
type Request struct {
	answer  chan Response
	message string
}

// NewRequest initializes an instance of a Request
func NewRequest(message string) Request {
	request := Request{}

	request.message = message

	request.answer = make(chan Response)

	return request
}

// GetAnswerChannel returns the channel back to the issuer of the request
func (self *Request) GetAnswerChannel() chan Response {
	return self.answer
}

// GetMessage returns a string containing the request
func (self *Request) GetMessage() string {
	return self.message
}

// SetAnswerChannel defines the answer channel
func (self *Request) SetAnswerChannel(answer chan Response) {
	self.answer = answer
}

// BaseResponse is a default struct that satisfies the Response interface
type BaseResponse struct {
	value interface{}
}

func (self BaseResponse) SetValue(value interface{}) {
	self.value = value
}

func (self BaseResponse) GetValue() interface{} {
	return self.value
}
