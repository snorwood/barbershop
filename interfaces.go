package barbershop

// Response is used to communicate between from a target to the source of a request
type Response interface {
	GetValue() interface{}
}

type Agent interface {
	GetReceiveRequestChan() chan Request
	Start()
}
