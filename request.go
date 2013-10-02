package barbershop

// Request is used to make requests from participants to the manager
type Request struct {
	answer  chan<- Subscriber
	message string
	target  string
}

// NewRequest initializes an instance of a Request
func NewRequest(target, message string) Request {
	request := Request{}
	request.target = target
	request.message = message
	request.answer = make(chan Subscriber)

	return request
}

// GetAnswerChannel returns the channel back to the issuer of the request
func (self *Request) GetAnswerChannel() chan<- Subscriber {
	return self.answer
}

// GetMessage returns a string containing the request
func (self *Request) GetMessage() string {
	return self.message
}

func (self *Request) GetTarget() string {
	return self.target
}

// SetAnswerChannel defines the answer channel
func (self *Request) SetAnswerChannel(answer chan<- Subscriber) {
	self.answer = answer
}

// BaseResponse is a default struct that satisfies the Response interface
type BaseResponse struct {
	value    interface{}
	positive bool
}

func (self BaseResponse) SetValue(value interface{}) {
	self.value = value
}

func (self BaseResponse) GetValue() interface{} {
	return self.value
}

func (self BaseResponse) SetPositive(positive bool) {
	self.positive = positive
}

func (self BaseResponse) GetPositive() interface{} {
	return self.positive
}
