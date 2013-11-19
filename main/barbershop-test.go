package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
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
	log       string
	busy      bool
	id        int
	Stop      chan IDGroup
	Start     chan int
	TimeSlept chan chan float32
	IsBusy    chan chan bool
	Finish    chan bool
}

func (self *Barber) GoLive() {
	var haircutTimer <-chan time.Time = nil
	finished := false
	customerID := 0
	setBusy := make(chan bool)
	for !finished {
		select {
		case <-haircutTimer:
			haircutTimer = nil
			go func(stop chan IDGroup, id IDGroup) {
				stop <- id
				setBusy <- false
			}(self.Stop, IDGroup{self.id, customerID})

		case customerID = <-self.Start:
			self.busy = true
			haircutTimer = time.After(time.Duration(rand.Int31n(6)+3) * time.Millisecond)

		case timeSlept := <-self.TimeSlept:
			timeSlept <- 0

		case isBusy := <-self.IsBusy:
			isBusy <- self.busy

		case self.busy = <-setBusy:

		case <-self.Finish:
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
	Finish    chan bool
}

func newBarber(id int, stop chan IDGroup) *BarberReader {
	start := make(chan int)
	timeSlept := make(chan chan float32)
	isBusy := make(chan chan bool)
	finish := make(chan bool)

	localBarber := BarberReader{
		Stop:      stop,
		Start:     start,
		TimeSlept: timeSlept,
		IsBusy:    isBusy,
		ID:        id,
		Finish:    finish,
	}
	barber := &Barber{
		id:        id,
		log:       "",
		busy:      false,
		Stop:      stop,
		Start:     start,
		TimeSlept: timeSlept,
		IsBusy:    isBusy,
		Finish:    finish,
	}

	go barber.GoLive()

	return &localBarber
}

type Customer struct {
	log        string
	id         int
	TimeWaited chan chan float32
	Message    chan string
}

func (self *Customer) GoLive() {
	for {
		select {
		case timeWaited := <-self.TimeWaited:
			timeWaited <- 0
		case message := <-self.Message:
			self.log += message
			fmt.Println(message)
		}
	}
}

type CustomerReader struct {
	TimeWaited chan chan float32
	Message    chan string
	ID         int
}

func newCustomer(id int) *CustomerReader {
	timeWaited := make(chan chan float32)
	message := make(chan string)

	localCustomer := CustomerReader{ID: id, TimeWaited: timeWaited, Message: message}
	customer := Customer{log: "", id: id, TimeWaited: timeWaited, Message: message}
	go customer.GoLive()

	return &localCustomer
}

var numBarbers int = 20

func main() {
	stop := make(chan IDGroup)
	barbers := make([]*BarberReader, 3)
	for id := 1; id <= 3; id++ {
		barbers[id-1] = newBarber(id, stop)
	}

	customers := make([]*CustomerReader, 10)
	customerCount := 0
	customersEntered := 0
	customersServed := 0

	for customersEntered < numBarbers || customerCount > 0 || !allBarbersFinished(barbers) {
		// newCustomerTimer := time.After(time.Duration(rand.Int31n(1)+1) * time.Second)
		newCustomerTimer := time.After(time.Millisecond)
		if customersEntered >= numBarbers {
			newCustomerTimer = nil
		}

		select {
		case <-newCustomerTimer:
			customersEntered += 1
			customer := newCustomer(customersEntered)
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
						break
					}
				}
			}

			if customerCount >= 10 && !foundBarber {
				fmt.Println("Customer", customer.ID, "was turned away.")
			} else if !foundBarber {
				customers[customerCount] = customer
				fmt.Println("Customer", customer.ID, "lined up.")
				customerCount += 1
			}
		case id := <-stop:
			customersServed += 1
			fmt.Println("Barber", id.BID, "stopped cutting customer", id.CID, "'s hair.")
			c := make(chan bool)
			barbers[id.BID-1].IsBusy <- c
			busy := <-c
			if !busy && customerCount > 0 {
				customer, _ := RemoveCustomer(customers, 0)
				fmt.Println("Barber", id.BID, "started cutting customer", customer.ID, "'s hair.")
				barbers[id.BID-1].Start <- customer.ID
				customerCount -= 1
			}
		}
	}

	for _, barber := range barbers {
		barber.Finish <- true
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
