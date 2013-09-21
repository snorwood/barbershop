package barbershop

type Manager struct {
	receiveRequestChan chan Request
}

func NewManager() *Manager {
	manager := new(Manager)
	manager.receiveRequestChan = make(chan Request, 100)
	go Start()

	return manager
}

func (self *Manager) GetRequestChan() chan Request {
	return self.receiveRequestChan
}

func (self *Manager) Start() {
	for {
		request := <-self.receiveRequestChan
		request.GetAnswerChannel() <- request.GetMessage()
	}
}

func (self *Manager) PrintSup() {
	print("SUP")
}
