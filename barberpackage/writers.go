package barbershop

import (
	"fmt"
	"io"
	"time"
)

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

type MessageRelay struct {
	incoming chan string
	outgoing chan string
	Dump     chan chan string
	dump     SWriter
}

func NewMessageRelay(incoming chan string, timeout time.Duration) MessageRelay {
	newRelay := MessageRelay{
		incoming: make(chan string),
		outgoing: make(chan string),
		Dump:     make(chan chan string),
		dump:     new(SWriter),
	}

	go func() {
		var timeoutChan chan bool = nil
		var internalOutgoing = nil
		var message string = nil
		incomingOpen := true

		for incomingOpen {
			select {
			case message, incomingOpen = <-newRelay.incoming:
				timeoutChan = time.After(timeout)
				internalOutgoing = newRelay.outgoing
			case internalOutgoing <- message:
				message = nil
				internalOutgoing = nil
				timeoutChan = nil
			case <-timeoutChan:
				fmt.Fprintln(newRelay.dump, message)
				message = nil
				internalOutgoing = nil
				timeoutChan = nil
			case sendDump := <-newRelay.Dump:
				sendDump <- newRelay.dump
			}
		}
	}()

	return newRelay
}

func (self MessageRelay) Outgoing() chan string {
	return self.outgoing
}

func (self MessageRelay) Write(p []byte) (int, error) {
	self.incoming <- string(p)
	return len(p), nil
}
