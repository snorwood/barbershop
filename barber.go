package barbershop

// Barber is a man who works at a barbershop gets paid to shave your face
type Barber struct {
	recieveRequestChan chan Request
}

// NewBarber initiallizes an instance of a Barber
func NewBarber(recieveRequestChan chan Request) *Barber {
	barber := new(Barber)
	barber.recieveRequestChan = recieveRequestChan

	return barber
}

// Start launches the barber's independent routine (go Start())
func (self *Barber) Start() {
	request := <-self.recieveRequestChan
	response := BaseResponse{value: request.message, positive: true}
	request.GetAnswerChannel() <- response
}

// SetRecieveRequestChan defines the channel the barber recieves requests on
func (self *Barber) SetRecieveRequestChan(recieveRequestChan chan Request) {
	self.recieveRequestChan = recieveRequestChan
}
