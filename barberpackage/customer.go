package barbershop

// Customer is a struct designed to live in its own routine and acts like a customer.
type Customer struct {
	log        *SWriter
	id         int
	TimeWaited chan chan time.Duration
	Log        chan chan string
	Stop       <-chan bool
	Start      chan int
	Enter      chan bool
	LineUp     chan int
	Kill       chan bool
}

// GoLive is the function to run in the customers own routine
func (self *Customer) GoLive() {

	// Define common variables
	inLine := false
	timeBegin := time.Now()
	stepTimer := timeBegin

	finished := false
	for !finished {
		select {
		//  Triggers when asked for timewaited. Sends time waited back.
		case timeWaited := <-self.TimeWaited:
			timeWaited <- time.Now().Sub(stepTimer)

		// Triggers when asked for log. Sends log back.
		case logger := <-self.Log:
			logger <- self.log.content

		// Triggers when haircut ends
		case <-self.Stop:
			fmt.Fprint(self.log, "Haircut finished. That took ", int(time.Now().Sub(timeBegin)/time.Second), " seconds total.")

		// Triggers when haircut starts. Barber id is recieved.
		case id := <-self.Start:
			if inLine {
				fmt.Fprint(self.log, "Haircut started with barber ", id, ". Waited for ", int(time.Now().Sub(stepTimer)/time.Second), " seconds.")
			} else {
				fmt.Fprint(self.log, "Haircut started with barber ", id, ". No wait :D.")
				stepTimer = time.Now()
			}
			inLine = false

		// Triggers when entering the store
		case <-self.Enter:
			fmt.Fprint(self.log, "Entered the store at ", time.Now().Format("15:04:05 on Jan 2"), ". My ID is ", self.id, ".")

		// Triggers when sent into line. Position in line is recieved.
		case id := <-self.LineUp:
			if id > 0 {
				inLine = true
				fmt.Fprint(self.log, "Got in line at position ", id, ".")
			} else {
				fmt.Fprint(self.log, "I was turned away.")
			}

		// Triggers when disposing of customer
		case <-self.Kill:
			finished = true
		}
	}
}

// CustomerReader is used for communicating with a customer.
type CustomerReader struct {
	TimeWaited chan chan time.Duration
	ID         int
	Log        chan chan string
	Stop       chan bool
	Start      chan int
	Enter      chan bool
	LineUp     chan int
	Kill       chan bool
}

// NewCustomer creates a customer and a customer reader. It starts the customer and returns the reader
func NewCustomer(id int) *CustomerReader {

	// Initialize shared channels
	timeWaited := make(chan chan time.Duration)
	log := make(chan chan string)
	stop := make(chan bool)
	start := make(chan int)
	enter := make(chan bool)
	lineUp := make(chan int)
	kill := make(chan bool)

	localCustomer := CustomerReader{
		ID:         id,
		TimeWaited: timeWaited,
		Log:        log,
		Stop:       stop,
		Start:      start,
		Enter:      enter,
		LineUp:     lineUp,
		Kill:       kill,
	}
	customer := Customer{
		log:        new(SWriter),
		id:         id,
		TimeWaited: timeWaited,
		Log:        log,
		Stop:       stop,
		Start:      start,
		Enter:      enter,
		LineUp:     lineUp,
		Kill:       kill,
	}

	// Start customer
	go customer.GoLive()

	return &localCustomer
}

// BestCustomer finds the customer who has been waiting the longest
func BestCustomer(customers []*CustomerReader) (*CustomerReader, error) {
	bestTime := time.Duration(0)
	var bestCustomer *CustomerReader = nil

	t := make(chan time.Duration)
	for _, customer := range customers {
		customer.TimeWaited <- t
		newTime := <-t
		if newTime > bestTime {
			bestTime = newTime
			bestCustomer = customer
		}
	}

	if len(customers) == 0 {
		return nil, fmt.Errorf("No customers in list")
	}

	if bestCustomer != nil {
		return bestCustomer, nil
	}

	return nil, fmt.Errorf("All customers are busy")
}

// RemoveCustomer removes the passed customer from the passed array
func RemoveCustomer(customers []*CustomerReader, customer *CustomerReader) ([]*CustomerReader, error) {
	index := -1

	for i := 0; i < len(customers); i++ {
		if customers[i] == customer {
			index = i
		}
	}

	if index == -1 {
		return customers, fmt.Errorf("Customer does not exist")
	}

	for visitor := index; visitor < len(customers)-1; visitor++ {
		customers[visitor] = customers[visitor+1]
	}

	if len(customers) > 0 {
		customers = customers[:len(customers)-1]
	}

	return customers, nil
}
