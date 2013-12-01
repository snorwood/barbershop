package barbershop

import (
	"io"
	"time"
)

// IDGroup contains ids of two communicators
type IDGroup struct {
	BID int
	CID int
}

type BarberShop struct {
	Output chan string
}

// Define constants
const (
	numCustomers int = 10
	numBarbers   int = 3
	variance     int = 1
	haircutBase  int = 4
	customerBase int = 1
)

func (shop *BarberShop) Start() chan string {
	w := NewMessageRelay(5 * time.Second)
	output := w.Outgoing()
	go shop.simulator(w)
	return output
}

func (shop *BarberShop) simulator(writer io.Writer) {
	// Initialize barbers and channel where they tell you they are done cutting hair
	stop := make(chan IDGroup)
	barbers := make([]*BarberReader, numBarbers)
	for id := 1; id <= numBarbers; id++ {
		barbers[id-1] = NewBarber(id, stop)
	}

	// Initialize customer helpers
	allCustomers := make([]*CustomerReader, numCustomers)
	customers := make([]*CustomerReader, 0, 10)
	customersEntered := 0
	newCustomerTimer := time.After(time.Duration(int(rand.Int31n(int32(variance)))+customerBase) * time.Second)

	// Loop until all customers have spawned and had their haircut / been turned away
	for customersEntered < numCustomers || len(customers) > 0 || !allBarbersFinished(barbers) {
		if customersEntered >= numCustomers {
			newCustomerTimer = nil
		}

		select {
		// Triggers when a new customer has spawned
		case <-newCustomerTimer:
			// Reset the timer
			newCustomerTimer = time.After(time.Duration(int(rand.Int31n(int32(variance)))+customerBase) * time.Second)

			// Enter customer
			customersEntered += 1
			customer := NewCustomer(customersEntered)
			allCustomers[customersEntered-1] = customer
			customer.Enter <- true
			fmt.Println("Customer", customer.ID, "entered.")

			// Check if barber is available
			foundBarber := false
			if len(customers) == 0 {
				barber, err := BestBarber(barbers)
				if err == nil {
					foundBarber = true
					fmt.Printf("Barber %d started cutting customer %d's hair.\n", barber.ID, customer.ID)
					barber.Start <- customer.ID
					customer.Start <- barber.ID
					break
				}
			}

			// Check if spot in line is available
			if len(customers) >= 10 && !foundBarber {
				customer.LineUp <- -1
				fmt.Println("Customer", customer.ID, "was turned away.")
			} else if !foundBarber {
				customers = customers[:len(customers)+1]
				customers[len(customers)-1] = customer
				fmt.Println("Customer", customer.ID, "lined up.")
				customer.LineUp <- len(customers)
			}

		// Triggers when a barber finishes
		case id := <-stop:
			// Stop the customer
			allCustomers[id.CID-1].Stop <- true
			fmt.Printf("Barber %d stopped cutting customer %d's hair.\n", id.BID, id.CID)

			// Search to see if there is a next customer in line
			if len(customers) > 0 {
				customer, err := BestCustomer(customers)
				if err == nil {
					updatedCustomers, _ := RemoveCustomer(customers, customer)
					customers = updatedCustomers
					fmt.Printf("Barber %d started cutting customer %d's hair.\n", id.BID, customer.ID)
					barbers[id.BID-1].Start <- customer.ID
					customer.Start <- id.BID
				}
			}
		}
	}

	// Tell barbers they are done cutting hair
	for _, barber := range barbers {
		barber.End <- true
	}

	// View logs
	// char := ""
	// reader := bufio.NewReader(os.Stdin)
	// for char != "q" {
	// 	fmt.Print("\nLog: Enter b for barbers, c for customers or q to quit: ")
	// 	input, _ := reader.ReadString('\n')
	// 	char = string([]byte(input)[0])
	// 	if char == "b" {
	// 		fmt.Print("\nEnter id of barber to view: ")
	// 		input, _ := reader.ReadString('\n')
	// 		char = string([]byte(input)[:len(input)-2])
	// 		i, err := strconv.Atoi(char)
	// 		if err == nil && i > 0 && i <= len(barbers) {
	// 			history := make(chan string)
	// 			barbers[i-1].Log <- history
	// 			fmt.Println("\n" + <-history)
	// 		}
	// 	} else if char == "c" {
	// 		fmt.Print("\nEnter id of customer to view: ")
	// 		input, _ := reader.ReadString('\n')
	// 		char = string([]byte(input)[:len(input)-2])
	// 		i, err := strconv.Atoi(char)
	// 		if err == nil && i > 0 && i <= len(allCustomers) {
	// 			history := make(chan string)
	// 			allCustomers[i-1].Log <- history
	// 			fmt.Println("\n" + <-history)
	// 		}
	// 	}
	// }

	// Dispose of customers
	for _, customer := range allCustomers {
		customer.Kill <- true
	}

	// Dispose of barbers
	for _, barber := range barbers {
		barber.Kill <- true
	}
}
