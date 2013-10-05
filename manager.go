package barbershop

import (
	"fmt"
	"math/rand"
	"time"
)

// Manager runs the communications in the barbershop
type Manager struct {
	customers          *CustomerQueue
	receiveRequestChan chan Request
	sendRequestChan    chan Request
	barbers            []chan Request
}

// NewManager initializes an instance of a manager
func NewManager() *Manager {
	manager := new(Manager)
	manager.receiveRequestChan = make(chan Request, 100)
	manager.sendRequestChan = make(chan Request)
	manager.customers = NewCustomerQueue(10)
	b := make([]Agent, 3)

	for i, _ := range b {
		b[i] = new(Barber)
	}

	manager.AddAgents(b)
	b = nil

	return manager
}

func (self *Manager) AddAgents(agents []Agent) error {
	err := false
	for _, agent := range agents {
		ch := make(chan Request)
		switch agent := agent.(type) {
		case *Barber:
			self.barbers = append(self.barbers, ch)
			agent.SetRecieveRequestChan(ch)
			go agent.Start()
		default:
			err = true
		}
	}

	if err {
		return nil
	} else {
		return fmt.Errorf("One or more invalid types were passed and were not added")
	}
}

func (self *Manager) GetRequestChan() chan Request {
	return self.receiveRequestChan
}

func (self *Manager) GetSendRequestChan() chan Request {
	return self.sendRequestChan
}

// Start initializes the managers separate routine (go Start())
func (self *Manager) Start() {
	customerReceive := make(chan *Customer)
	customerReceiveBackup := customerReceive
	go customerGenerator(20, customerReceive)

	for {
		if self.customers.Full() {
			customerReceive = nil
			break
		} else {
			customerReceive = customerReceiveBackup
		}

		select {
		case request := <-self.receiveRequestChan:
			self.redirect(request)
		case customer := <-customerReceive:
			self.customers.Enqueue(customer)
		}
	}
}

func (self *Manager) redirect(request Request) {
	if request.GetTarget() == "Barber" {
		for _, barber := range self.barbers {
			barber <- request
		}
	}
}

func customerGenerator(numberOfCustomers int, send chan *Customer) {
	for i := 1; i <= numberOfCustomers; i++ {
		rand.Seed(time.Now().UnixNano())
		customer := &Customer{}
		<-time.After(time.Duration(rand.Int31n(4)) * time.Second)
		select {
		case send <- customer:
			fmt.Printf("Customer %d was seated\n", i)
		default:
		}
	}
}
