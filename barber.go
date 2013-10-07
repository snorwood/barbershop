package barbershop

// Barber is a man who works at a barbershop gets paid to shave your face
type Barber struct {
	recieveRequestChan chan Request
}

// NewBarber initiallizes an instance of a Barber
func NewBarber() *Barber {
	barber := new(Barber)
	barber.recieveRequestChan = make(chan Request)

	return barber
}

// Start launches the barber's independent routine (go Start())
func (self *Barber) Start() {
	request := <-self.recieveRequestChan
	subscriber := NewSubscriber()
	SendSubscriber(subscriber, request.GetAnswerChannel())
	select {
	case message := <-subscriber.Receive:
		print(message)
	case <-subscriber.StopReceiving:
	}
}

func (self *Barber) GetReceiveRequestChan() chan Request {
	return self.recieveRequestChan
}

// SetRecieveRequestChan defines the channel the barber recieves requests on
func (self *Barber) SetRecieveRequestChan(recieveRequestChan chan Request) {
	self.recieveRequestChan = recieveRequestChan
}
