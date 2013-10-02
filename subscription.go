package barbershop

type Subscriber struct {
	Send          chan string
	Receive       chan string
	StopReceiving chan bool
	StopSending   chan bool
}

func NewSubscriber() Subscriber {
	subscriber := new(Subscriber)
	subscriber.Send = make(chan string)
	subscriber.Receive = make(chan string)
	subscriber.StopReceiving = make(chan bool)
	subscriber.StopSending = make(chan bool)
	
	return *subscriber
}

func SendSubscriber(subscriber Subscriber, send chan<- Subscriber) {
	newSubscriber := Subscriber{
		Send:          subscriber.Receive,
		Receive:       subscriber.Send,
		StopReceiving: subscriber.StopSending,
		StopSending:   subscriber.StopReceiving,
	}

	send <- newSubscriber
}
