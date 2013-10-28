package main

import (
	"math/rand"
	"time"
)

type person struct{
	bool Busy
	string Text 
}

func spawnCustomer(entrance chan person) {
	for i := 0; i < 30; i++ {
		random := rand.Int(10) * time.Second
		<- time.After(random)

		if (cap(line) < len(line))
		{
			person := Person{}
			entrance <- person
		}
	}
}

func main() {
	barbers = make([]person, 3)
	line = make([]person, 0, 10)

	entrance := make(chan person)
	doneCutting := make(chan chan bool)

	go spawnCustomer(entrance)

	select {
	case person := <-entrance:
		for _, barber := range(barbers)
		{
			if (!barber.Busy)
			{
				barber.Text = "c"
				barber.Busy = true
				go func f() {
					<- time.After((rand.int(3) + 3) * time.Second)
					c := make(chan bool)
					doneCutting <- c
					sleep := !(<- c)
					if (sleep)
					{
						barber.Text = "z"
						barber.Busy = false
					} else {
						barber.Busy = true
						barber.Text = "c"
						f()
					}
				}()
			} else {
				
			}
		}
	case c := <-doneCutting:
		if (cap(line) > 0)
		{
			c <- true
			barber.Busy = true
			barber.Text = "c"
		}
		else
		{
			c <- false
		}
	}

}
