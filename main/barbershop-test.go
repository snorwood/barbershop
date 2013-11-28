package barbershop

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// IDGroup contains ids of two communicators
type IDGroup struct {
	BID int
	CID int
}

// MyWriter writes every input on a new line with line numbers to multiple outputs
type MyWriter struct {
	count   int
	outputs []io.Writer
}

// SWriter writes every input to a string with line numbers
type SWriter struct {
	count   int
	content string
}

// Write input to persistant string
func (self *SWriter) Write(p []byte) (int, error) {
	self.count += 1
	self.content += fmt.Sprintf("%d.\t %s\n", self.count, string(p))

	return len(p), nil
}

// Write input to contained output
func (self *MyWriter) Write(p []byte) (int, error) {
	self.count += 1
	s := fmt.Sprintf("%d.\t %s\n", self.count, string(p))

	for _, output := range self.outputs {
		fmt.Fprint(output, s)
	}
	return len(p), nil
}

// NewMyWriter creates a new MyWriter struct
func NewMyWriter(outputs ...io.Writer) *MyWriter {
	writer := MyWriter{outputs: outputs}
	return &writer
}

// Barber is a struct designed to run in its own routine and act like a barber
type Barber struct {
	log       *SWriter
	busy      bool
	id        int
	Stop      chan IDGroup
	Start     chan int
	TimeSlept chan chan time.Duration
	IsBusy    chan chan bool
	Log       chan chan string
	End       chan bool
	Kill      chan bool
}

// GoLive is the function that runs in the barber's routine
func (self *Barber) GoLive() {
	// Sends when haircut is over
	var haircutTimer <-chan time.Time = nil

	// Initializes variables to be used
	customerID := 0
	setBusy := make(chan bool)
	timeBegin := time.Now()

	fmt.Fprint(self.log, "Became sentient. My ID is ", self.id)
	finished := false
	for !finished {
		select {
		// Triggers when haircut is over
		case <-haircutTimer:
			fmt.Fprint(self.log, "Finished cutting customer ", customerID, "'s hair. Cut took ", int64(time.Now().Sub(timeBegin)/time.Second), " seconds.")
			timeBegin = time.Now()
			haircutTimer = nil

			// sending inline can cause stalling
			go func(stop chan IDGroup, id IDGroup) {
				stop <- id
				setBusy <- false
			}(self.Stop, IDGroup{self.id, customerID})

		// Triggers at the start of a haircut
		case customerID = <-self.Start:
			fmt.Fprint(self.log, "Started cutting customer ", customerID, "'s hair. Slept for ", int(time.Now().Sub(timeBegin)/time.Second), " seconds.")
			timeBegin = time.Now()
			self.busy = true
			haircutTimer = time.After(time.Duration(int(rand.Int31n(int32(variance)))+haircutBase) * time.Second)

		// Triggers on request for timeSlept. Sends time slept back.
		case timeSlept := <-self.TimeSlept:
			timeSlept <- time.Now().Sub(timeBegin)

		// Triggers on request for busy. Sends busy back.
		case isBusy := <-self.IsBusy:
			isBusy <- self.busy

		// Triggers when setting busy
		case self.busy = <-setBusy:

		// Triggers on request for log. Sends log back.
		case logger := <-self.Log:
			logger <- self.log.content

		// Triggers when simulation is over.
		case <-self.End:
			fmt.Fprint(self.log, "Done for the day. Phew!")

		// Triggers to dispose of the barber
		case <-self.Kill:
			finished = true
		}
	}
}

// BarberReader is a struct for communicating with a barber. Going to combine into just a barber soon.
type BarberReader struct {
	Stop      chan IDGroup
	Start     chan int
	TimeSlept chan chan time.Duration
	IsBusy    chan chan bool
	ID        int
	Log       chan chan string
	End       chan bool
	Kill      chan bool
}

// NewBarber creates a barber and a barber reader. It starts the barber in its own routine and returns the reader.
func NewBarber(id int, stop chan IDGroup) *BarberReader {

	// Initialize shared channels
	start := make(chan int)
	timeSlept := make(chan chan time.Duration)
	isBusy := make(chan chan bool)
	log := make(chan chan string)
	end := make(chan bool)
	kill := make(chan bool)

	localBarber := &BarberReader{
		Stop:      stop,
		Start:     start,
		TimeSlept: timeSlept,
		IsBusy:    isBusy,
		ID:        id,
		Log:       log,
		End:       end,
		Kill:      kill,
	}

	barber := &Barber{
		id:        id,
		log:       new(SWriter),
		busy:      false,
		Stop:      stop,
		Start:     start,
		TimeSlept: timeSlept,
		IsBusy:    isBusy,
		Log:       log,
		End:       end,
		Kill:      kill,
	}

	// Start barber go routine
	go barber.GoLive()

	return localBarber
}

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

// Define constants
const (
	numCustomers int = 10
	numBarbers   int = 3
	variance     int = 1
	haircutBase  int = 4
	customerBase int = 1
)

func (shop *BarberShop) Start() {
	go shop.simulator()
}

func (shop *BarberShop) simulator() {
	// Initialize barbers and channel where they tell you they are done cutting hair
	stop := make(chan IDGroup)
	barbers := make([]*BarberReader, numBarbers)
	for id := 1; id <= numBarbers; id++ {
		barbers[id-1] = NewBarber(id, stop)
	}

	// Initialize customer helpers
	allCustomers := make([]*CustomerReader, numCustomers)
	customers := make([]*CustomerReader, 0, 10)
	customersEntered := 0
	newCustomerTimer := time.After(time.Duration(int(rand.Int31n(int32(variance)))+customerBase) * time.Second)

	// Loop until all customers have spawned and had their haircut / been turned away
	for customersEntered < numCustomers || len(customers) > 0 || !allBarbersFinished(barbers) {
		if customersEntered >= numCustomers {
			newCustomerTimer = nil
		}

		select {
		// Triggers when a new customer has spawned
		case <-newCustomerTimer:
			// Reset the timer
			newCustomerTimer = time.After(time.Duration(int(rand.Int31n(int32(variance)))+customerBase) * time.Second)

			// Enter customer
			customersEntered += 1
			customer := NewCustomer(customersEntered)
			allCustomers[customersEntered-1] = customer
			customer.Enter <- true
			fmt.Println("Customer", customer.ID, "entered.")

			// Check if barber is available
			foundBarber := false
			if len(customers) == 0 {
				barber, err := BestBarber(barbers)
				if err == nil {
					foundBarber = true
					fmt.Printf("Barber %d started cutting customer %d's hair.\n", barber.ID, customer.ID)
					barber.Start <- customer.ID
					customer.Start <- barber.ID
					break
				}
			}

			// Check if spot in line is available
			if len(customers) >= 10 && !foundBarber {
				customer.LineUp <- -1
				fmt.Println("Customer", customer.ID, "was turned away.")
			} else if !foundBarber {
				customers = customers[:len(customers)+1]
				customers[len(customers)-1] = customer
				fmt.Println("Customer", customer.ID, "lined up.")
				customer.LineUp <- len(customers)
			}

		// Triggers when a barber finishes
		case id := <-stop:
			// Stop the customer
			allCustomers[id.CID-1].Stop <- true
			fmt.Printf("Barber %d stopped cutting customer %d's hair.\n", id.BID, id.CID)

			// Search to see if there is a next customer in line
			if len(customers) > 0 {
				customer, err := BestCustomer(customers)
				if err == nil {
					updatedCustomers, _ := RemoveCustomer(customers, customer)
					customers = updatedCustomers
					fmt.Printf("Barber %d started cutting customer %d's hair.\n", id.BID, customer.ID)
					barbers[id.BID-1].Start <- customer.ID
					customer.Start <- id.BID
				}
			}
		}
	}

	// Tell barbers they are done cutting hair
	for _, barber := range barbers {
		barber.End <- true
	}

	// View logs
	char := ""
	reader := bufio.NewReader(os.Stdin)
	for char != "q" {
		fmt.Print("\nLog: Enter b for barbers, c for customers or q to quit: ")
		input, _ := reader.ReadString('\n')
		char = string([]byte(input)[0])
		if char == "b" {
			fmt.Print("\nEnter id of barber to view: ")
			input, _ := reader.ReadString('\n')
			char = string([]byte(input)[:len(input)-2])
			i, err := strconv.Atoi(char)
			if err == nil && i > 0 && i <= len(barbers) {
				history := make(chan string)
				barbers[i-1].Log <- history
				fmt.Println("\n" + <-history)
			}
		} else if char == "c" {
			fmt.Print("\nEnter id of customer to view: ")
			input, _ := reader.ReadString('\n')
			char = string([]byte(input)[:len(input)-2])
			i, err := strconv.Atoi(char)
			if err == nil && i > 0 && i <= len(allCustomers) {
				history := make(chan string)
				allCustomers[i-1].Log <- history
				fmt.Println("\n" + <-history)
			}
		}
	}

	// Dispose of customers
	for _, customer := range allCustomers {
		customer.Kill <- true
	}

	// Dispose of barbers
	for _, barber := range barbers {
		barber.Kill <- true
	}
}

// BestBarber finds the barber who has been waiting the longest
func BestBarber(barbers []*BarberReader) (*BarberReader, error) {
	bestTime := time.Duration(0)
	var bestBarber *BarberReader = nil

	t := make(chan time.Duration)
	b := make(chan bool)
	for _, barber := range barbers {
		barber.IsBusy <- b
		if !<-b {
			barber.TimeSlept <- t
			newTime := <-t
			if newTime > bestTime {
				bestTime = newTime
				bestBarber = barber
			}
		}
	}

	if len(barbers) == 0 {
		return nil, fmt.Errorf("No barbers in list")
	}

	if bestBarber != nil {
		return bestBarber, nil
	}

	return nil, fmt.Errorf("All barbers are busy")
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

// allBarbersFinished checks if all of the barbers are free (not busy)
func allBarbersFinished(barbers []*BarberReader) bool {
	for _, barber := range barbers {
		c := make(chan bool)
		barber.IsBusy <- c
		busy := <-c
		if busy {
			return false
		}
	}

	return true
}
