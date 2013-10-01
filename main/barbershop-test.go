package main

import (
	"github.com/snorwood/barbershop"
)

func main() {
	manager := barbershop.NewManager()
	request := barbershop.NewRequest("Barber", "Hello")
	customer := barbershop.NewCustomer(manager.GetRequestChan())
	subscription := <-customer.SendRequest(request)
	if str, ok := response.GetValue().(string); ok {
		print(str)
	} else {
		print(":( no string for you")
	}
}
