package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type Person struct {
	Busy bool
	ID   int
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
	for i := 1; i <= 30; i++ {
		random := time.Duration(rand.Int31n(2)) * time.Second
		<-time.After(random)
		person := &Person{ID: i}
		entrance <- person
		fmt.Fprint(output, "Added new Customer ", i)
	}

	done <- true
}

func main() {
	os.Create("output.txt")
	log, _ := os.OpenFile("output.txt", os.O_APPEND, 0666)
	defer log.Close() //we'll close this file as we leave scope, no matter what

	output = NewMyWriter(os.Stdout, log)

	barbers := make([]*Person, 3)
	line := make([]*Person, 0, 10)

	entrance := make(chan *Person)
	doneCutting := make(chan chan bool)
	done := make(chan bool)
	go spawnCustomer(entrance, done)
	running := true

	for i := range barbers {
		barbers[i] = &Person{Busy: false, ID: i}
	}

	for running {
		select {
		case person := <-entrance:
			spotAvailable := false
			for _, barber := range barbers {
				if !barber.Busy {
					spotAvailable = true
					var f func()
					f = func() {
						barber.Busy = true
						fmt.Fprint(output, "Barber ", barber.ID, " is gettin' to work on Customer ", person.ID)
						<-time.After(time.Duration(rand.Int31n(3)+39) * time.Second)
						c := make(chan bool)
						doneCutting <- c
						fmt.Fprint(output, "Barber ", barber.ID, " finished haircut on Customer ", person.ID)
						sleep := !(<-c)
						if sleep {
							fmt.Fprint(output, "Barber ", barber.ID, " is snoozin'")
							barber.Busy = false
						} else {
							f()
						}
					}
					go f()
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
				line = append(line[:0], line[1:]...)
				c <- true
			} else {
				c <- false
			}
		case <-done:
			running = false
		}

	}
}
