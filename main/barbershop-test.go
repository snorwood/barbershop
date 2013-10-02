package main

import (
	"github.com/snorwood/barbershop"
)

func main() {
	manager := barbershop.NewManager()
	request := barbershop.NewRequest("Barber", "Hello")
	customer := barbershop.NewCustomer(manager.GetRequestChan())
	subscription := <-customer.SendRequest(request)
	select {
	case event := <-subscription.Receive:
		println(event)
	default:
		println("not successful")
	}
	subscription.StopReceiving <- true
	close(subscription.StopReceiving)
	<-subscription.StopSending
	close(subscription.Send)
	println("Cleaned up properly")
}
