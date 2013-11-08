package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type Person struct {
	busy bool
	ID   int
}

type Barber struct {
	busy       bool
	ReturnBusy chan chan bool
	SetBusy    chan bool

	ID int
}

func (self *Barber) DataStream() {
	self.ReturnBusy = make(chan chan bool)
	self.SetBusy = make(chan bool)
	for {
		select {
		case returnBusy := <-self.ReturnBusy:
			returnBusy <- self.busy
		case busy := <-self.SetBusy:
			self.busy = busy
		}
	}
}

func (self *Barber) GetBusy() bool {
	c := make(chan bool)
	self.ReturnBusy <- c
	return <-c
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

func spawnCustomer(entrance chan *Person, done chan bool) {
	for i := 1; i <= 10; i++ {
		random := time.Duration(rand.Int31n(2)+1) * time.Second
		<-time.After(random)
		person := &Person{ID: i}
		fmt.Fprint(output, "Added new Customer ", i)
		entrance <- person
	}

	done <- true
}

func checkBusy(barbers []*Barber) bool {
	for _, barber := range barbers {
		if barber.GetBusy() {
			return true
		}
	}
	return false
}

func HairCut(barber *Barber, person *Person, doneCutting chan chan *Person) {
	barber.SetBusy <- true
	fmt.Fprint(output, "Barber ", barber.ID, " is gettin' to work on Customer ", person.ID)
	<-time.After(time.Duration(rand.Int31n(3)+3) * time.Second)
	fmt.Fprint(output, "Barber ", barber.ID, " finished haircut on Customer ", person.ID)
	barber.SetBusy <- false
	c := make(chan *Person)
	doneCutting <- c
	person = (<-c)
	if person == nil {
		fmt.Fprint(output, "Barber ", barber.ID, " is snoozin'")
	} else {
		go HairCut(barber, person, doneCutting)
	}
}

func main() {
	os.Create("output.txt")
	log, _ := os.OpenFile("output.txt", os.O_APPEND, 0666)
	defer log.Close() //we'll close this file as we leave scope, no matter what

	output = NewMyWriter(os.Stdout, log)

	barbers := make([]*Barber, 3)
	line := make([]*Person, 0, 10)

	entrance := make(chan *Person, 10)
	doneCutting := make(chan chan *Person)
	done := make(chan bool)
	go spawnCustomer(entrance, done)
	noMoreCustomers := false

	for i := range barbers {
		barbers[i] = &Barber{busy: false, ID: i + 1}
		go barbers[i].DataStream()
	}

	for !noMoreCustomers || checkBusy(barbers) {
		select {
		case person := <-entrance:
			spotAvailable := false
			for _, barber := range barbers {
				if !barber.GetBusy() {
					spotAvailable = true
					go HairCut(barber, person, doneCutting)
					break
				}
			}

			if !spotAvailable {
				if len(line) < cap(line) {
					fmt.Fprint(output, "Customer ", person.ID, " got in line")
					line = append(line, person)
				} else {
					fmt.Fprint(output, "Customer ", person.ID, " has been turned away")
				}
			}

		case c := <-doneCutting:
			if len(line) > 0 {
				c <- line[0]
				line = append(line[:0], line[1:]...)
			} else {
				c <- nil
			}

		case <-done:
			noMoreCustomers = true
		}
	}
}
