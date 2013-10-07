package barbershop

import (
	"fmt"
	// "time"
)

// Manager runs the communications in the barbershop
type Manager struct {
	waitingRoom        *GeneralQueue
	receiveRequestChan chan Request
	sendRequestChan    chan Request
	barbers            []chan Request
	customers          []chan Request
}

// NewManager initializes an instance of a manager
func NewManager() *Manager {
	manager := new(Manager)
	manager.receiveRequestChan = make(chan Request, 100)
	manager.sendRequestChan = make(chan Request)
	manager.waitingRoom = NewGeneralQueue(10)
	b := make([]Agent, 3)

	for i, _ := range b {
		b[i] = new(Barber)
	}

	manager.AddAgents(b)
	b = nil

	return manager
}

func (self *Manager) AddAgent(agent Agent) error {
	ch := make(chan Request)
	switch agent := agent.(type) {
	case *Barber:
		self.barbers = append(self.barbers, ch)
		go agent.Start()
	case *Customer:
		self.customers = append(self.customers)
		go agent.Start()
	default:
		return fmt.Errorf("One or more invalid types were passed and were not added")
	}

	return nil
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
	}

	return fmt.Errorf("One or more invalid types were passed and were not added")
}

func (self *Manager) GetRequestChan() chan Request {
	return self.receiveRequestChan
}

func (self *Manager) GetSendRequestChan() chan Request {
	return self.sendRequestChan
}

// Start initializes the managers separate routine (go Start())
func (self *Manager) Start() {
	customerReceive := customerGenerator(1, self)
	customerReceiveBackup := customerReceive

	for {
		if self.waitingRoom.Full() {
			customerReceive = nil
			break
		} else {
			customerReceive = customerReceiveBackup
		}

		select {
		case request := <-self.receiveRequestChan:
			self.redirect(request)
		case customer := <-customerReceive:
			err := self.AddAgent(customer)
			if err == nil {
				self.waitingRoom.Enqueue(customer)
			}
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

func customerGenerator(numberOfCustomers int, manager *Manager) chan *Customer {
	send := make(chan *Customer)
	go func() {
		for i := 1; i <= numberOfCustomers; i++ {
			customer := NewCustomer(manager.GetRequestChan())

			send <- customer
		}
	}()
	return send
}
