package barbershop

// Manager runs the communications in the barbershop
type Manager struct {
	receiveRequestChan chan Request
	sendRequestChan    chan Request
}

// NewManager initializes an instance of a manager
func NewManager() *Manager {
	manager := new(Manager)
	manager.receiveRequestChan = make(chan Request, 100)
	manager.sendRequestChan = make(chan Request)
	go manager.Start()

	return manager
}

func (self *Manager) GetRequestChan() chan Request {
	return self.receiveRequestChan
}

func (self *Manager) GetSendRequestChan() chan Request {
	return self.sendRequestChan
}

// Start initializes the managers separate routine (go Start())
func (self *Manager) Start() {
	for {
		request := <-self.receiveRequestChan
		self.sendRequestChan <- request
	}
}
