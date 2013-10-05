package main

import (
	"github.com/snorwood/barbershop"
)

func main() {
	manager := barbershop.NewManager()
	manager.Start()
}
