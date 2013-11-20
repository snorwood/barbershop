package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type IDGroup struct {
	BID int
	CID int
}

type MyWriter struct {
	count   int
	outputs []io.Writer
}

type SWriter struct {
	count   int
	content string
}

func (self *SWriter) Write(p []byte) (int, error) {
	self.count += 1
	self.content += fmt.Sprintf("%d.\t %s\n", self.count, string(p))

	return len(p), nil
}

func (self *MyWriter) Write(p []byte) (int, error) {
	self.count += 1
	s := fmt.Sprintf("%d.\t %s\n", self.count, string(p))

	for _, output := range self.outputs {
		fmt.Fprint(output, s)
	}
	return len(p), nil
}

func NewMyWriter(outputs ...io.Writer) *MyWriter {
	writer := MyWriter{outputs: outputs}
	return &writer
}

var current = time.Now()
var log *os.File
var output io.Writer

type Barber struct {
	log       *SWriter
	busy      bool
	id        int
	Stop      chan IDGroup
	Start     chan int
	TimeSlept chan chan float32
	IsBusy    chan chan bool
	Log       chan chan string
	End       chan bool
	Kill      chan bool
}

func (self *Barber) GoLive() {
	var haircutTimer <-chan time.Time = nil
	finished := false
	customerID := 0
	setBusy := make(chan bool)
	fmt.Fprint(self.log, "Became sentient. My ID is ", self.id)
	timeBegin := time.Now()

	for !finished {
		select {
		case <-haircutTimer:
			fmt.Fprint(self.log, "Finished cutting customer ", customerID, "'s hair. Cut took ", int64(time.Now().Sub(timeBegin)/time.Second), " seconds.")
			timeBegin = time.Now()
			haircutTimer = nil
			go func(stop chan IDGroup, id IDGroup) {
				stop <- id
				setBusy <- false
			}(self.Stop, IDGroup{self.id, customerID})

		case customerID = <-self.Start:
			fmt.Fprint(self.log, "Started cutting customer ", customerID, "'s hair. Slept for ", int(time.Now().Sub(timeBegin)/time.Second), " seconds.")
			timeBegin = time.Now()
			self.busy = true
			haircutTimer = time.After(time.Duration(int(rand.Int31n(int32(haircutBase)))+variance) * time.Second)

		case timeSlept := <-self.TimeSlept:
			timeSlept <- 0

		case isBusy := <-self.IsBusy:
			isBusy <- self.busy

		case self.busy = <-setBusy:

		case logger := <-self.Log:
			logger <- self.log.content

		case <-self.End:
			fmt.Fprint(self.log, "Done for the day. Phew!")

		case <-self.Kill:
			finished = true
		}
	}
}

type BarberReader struct {
	Stop      chan IDGroup
	Start     chan int
	TimeSlept chan chan float32
	IsBusy    chan chan bool
	ID        int
	Log       chan chan string
	End       chan bool
	Kill      chan bool
}

func newBarber(id int, stop chan IDGroup) *BarberReader {
	start := make(chan int)
	timeSlept := make(chan chan float32)
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

	go barber.GoLive()

	return localBarber
}

type Customer struct {
	log        *SWriter
	id         int
	TimeWaited chan chan float32
	Message    chan string
	Log        chan chan string
	Stop       <-chan bool
	Start      chan int
	Enter      chan bool
	LineUp     chan int
	Kill       chan bool
}

func (self *Customer) GoLive() {
	finished := false
	inLine := false
	timeBegin := time.Now()
	stepTimer := timeBegin

	for !finished {
		select {
		case timeWaited := <-self.TimeWaited:
			timeWaited <- 0

		case message := <-self.Message:
			fmt.Fprintln(self.log, message)

		case logger := <-self.Log:
			logger <- self.log.content

		case <-self.Stop:
			fmt.Fprint(self.log, "Haircut finished. That took ", int(time.Now().Sub(timeBegin)/time.Second), " seconds total.")

		case id := <-self.Start:
			if inLine {
				fmt.Fprint(self.log, "Haircut started with barber ", id, ". No wait :D.")
			} else {
				fmt.Fprint(self.log, "Haircut started with barber ", id, ". Waited for ", int(time.Now().Sub(stepTimer)/time.Second), " seconds.")
				stepTimer = time.Now()
			}
			inLine = false

		case <-self.Enter:
			fmt.Fprint(self.log, "Entered the store at ", time.Now().Format("15:04:05 on Jan 2"), ". My ID is ", self.id, ".")

		case id := <-self.LineUp:
			if id > 0 {
				inLine = true
				fmt.Fprint(self.log, "Got in line at position ", id, ".")
			} else {
				fmt.Fprint(self.log, "I was turned away.")
			}

		case <-self.Kill:
			finished = true
		}
	}
}

type CustomerReader struct {
	TimeWaited chan chan float32
	Message    chan string
	ID         int
	Log        chan chan string
	Stop       chan bool
	Start      chan int
	Enter      chan bool
	LineUp     chan int
	Kill       chan bool
}

func newCustomer(id int) *CustomerReader {
	timeWaited := make(chan chan float32)
	message := make(chan string)
	log := make(chan chan string)
	stop := make(chan bool)
	start := make(chan int)
	enter := make(chan bool)
	lineUp := make(chan int)
	kill := make(chan bool)

	localCustomer := CustomerReader{
		ID:         id,
		TimeWaited: timeWaited,
		Message:    message,
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
		Message:    message,
		Log:        log,
		Stop:       stop,
		Start:      start,
		Enter:      enter,
		LineUp:     lineUp,
		Kill:       kill,
	}

	go customer.GoLive()

	return &localCustomer
}

var numCustomers int = 4
var numBarbers int = 3
var variance int = 1
var haircutBase int = 3
var customerBase int = 1

func main() {
	stop := make(chan IDGroup)
	barbers := make([]*BarberReader, 3)
	for id := 1; id <= numBarbers; id++ {
		barbers[id-1] = newBarber(id, stop)
	}

	allCustomers := make([]*CustomerReader, numCustomers)
	customers := make([]*CustomerReader, 10)
	customerCount := 0
	customersEntered := 0
	customersServed := 0

	for customersEntered < numCustomers || customerCount > 0 || !allBarbersFinished(barbers) {
		newCustomerTimer := time.After(time.Duration(int(rand.Int31n(int32(customerBase)))+variance) * time.Second)
		//newCustomerTimer := time.After(time.Millisecond)
		if customersEntered >= numCustomers {
			newCustomerTimer = nil
		}

		select {
		case <-newCustomerTimer:
			customersEntered += 1
			customer := newCustomer(customersEntered)
			allCustomers[customersEntered-1] = customer
			customer.Enter <- true
			fmt.Println("Customer", customer.ID, "entered.")
			foundBarber := false

			if customerCount == 0 {
				for _, barber := range barbers {
					c := make(chan bool)
					barber.IsBusy <- c
					busy := <-c
					if !busy {
						foundBarber = true
						fmt.Println("Barber", barber.ID, "started cutting customer", customer.ID, "'s hair.")
						barber.Start <- customer.ID
						customer.Start <- barber.ID
						break
					}
				}
			}

			if customerCount >= 10 && !foundBarber {
				customer.LineUp <- -1
				fmt.Println("Customer", customer.ID, "was turned away.")
			} else if !foundBarber {
				customers[customerCount] = customer
				fmt.Println("Customer", customer.ID, "lined up.")
				customerCount += 1
				customer.LineUp <- customerCount
			}

		case id := <-stop:
			customersServed += 1
			allCustomers[id.CID-1].Stop <- true
			fmt.Println("Barber", id.BID, "stopped cutting customer", id.CID, "'s hair.")
			c := make(chan bool)
			barbers[id.BID-1].IsBusy <- c
			busy := <-c
			if !busy && customerCount > 0 {
				customer, _ := RemoveCustomer(customers, 0)
				fmt.Println("Barber", id.BID, "started cutting customer", customer.ID, "'s hair.")
				barbers[id.BID-1].Start <- customer.ID
				customer.Start <- id.BID
				customerCount -= 1
			}
		}
	}

	for _, barber := range barbers {
		barber.End <- true
	}

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
		// s := make(chan string)
		// barbers[0].Log <- s
		// fmt.Println(<-s, "\n")
	}

	for _, customer := range allCustomers {
		customer.Kill <- true
	}

	for _, barber := range barbers {
		barber.Kill <- true
	}
}

func RemoveCustomer(customers []*CustomerReader, index int) (*CustomerReader, error) {
	if index >= len(customers) {
		return nil, fmt.Errorf("Array index out of bounds")
	}

	customer := customers[index]

	for visitor := index; visitor < len(customers)-1; visitor++ {
		customers[visitor] = customers[visitor+1]
	}

	return customer, nil
}

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
