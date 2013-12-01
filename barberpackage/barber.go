package barbershop

// Barber is a struct designed to run in its own routine and act like a barber
type Barber struct {
	log       *SWriter
	busy      bool
	id        int
	Stop      chan IDGroup
	Start     chan int
	TimeSlept chan chan time.Duration
	IsBusy    chan chan bool
	Log       chan chan string
	End       chan bool
	Kill      chan bool
}

// GoLive is the function that runs in the barber's routine
func (self *Barber) GoLive() {
	// Sends when haircut is over
	var haircutTimer <-chan time.Time = nil

	// Initializes variables to be used
	customerID := 0
	setBusy := make(chan bool)
	timeBegin := time.Now()

	fmt.Fprint(self.log, "Became sentient. My ID is ", self.id)
	finished := false
	for !finished {
		select {
		// Triggers when haircut is over
		case <-haircutTimer:
			fmt.Fprint(self.log, "Finished cutting customer ", customerID, "'s hair. Cut took ", int64(time.Now().Sub(timeBegin)/time.Second), " seconds.")
			timeBegin = time.Now()
			haircutTimer = nil

			// sending inline can cause stalling
			go func(stop chan IDGroup, id IDGroup) {
				stop <- id
				setBusy <- false
			}(self.Stop, IDGroup{self.id, customerID})

		// Triggers at the start of a haircut
		case customerID = <-self.Start:
			fmt.Fprint(self.log, "Started cutting customer ", customerID, "'s hair. Slept for ", int(time.Now().Sub(timeBegin)/time.Second), " seconds.")
			timeBegin = time.Now()
			self.busy = true
			haircutTimer = time.After(time.Duration(int(rand.Int31n(int32(variance)))+haircutBase) * time.Second)

		// Triggers on request for timeSlept. Sends time slept back.
		case timeSlept := <-self.TimeSlept:
			timeSlept <- time.Now().Sub(timeBegin)

		// Triggers on request for busy. Sends busy back.
		case isBusy := <-self.IsBusy:
			isBusy <- self.busy

		// Triggers when setting busy
		case self.busy = <-setBusy:

		// Triggers on request for log. Sends log back.
		case logger := <-self.Log:
			logger <- self.log.content

		// Triggers when simulation is over.
		case <-self.End:
			fmt.Fprint(self.log, "Done for the day. Phew!")

		// Triggers to dispose of the barber
		case <-self.Kill:
			finished = true
		}
	}
}

// BarberReader is a struct for communicating with a barber. Going to combine into just a barber soon.
type BarberReader struct {
	Stop      chan IDGroup
	Start     chan int
	TimeSlept chan chan time.Duration
	IsBusy    chan chan bool
	ID        int
	Log       chan chan string
	End       chan bool
	Kill      chan bool
}

// NewBarber creates a barber and a barber reader. It starts the barber in its own routine and returns the reader.
func NewBarber(id int, stop chan IDGroup) *BarberReader {

	// Initialize shared channels
	start := make(chan int)
	timeSlept := make(chan chan time.Duration)
	isBusy := make(chan chan bool)
	log := make(chan chan string)
	end := make(chan bool)
	kill := make(chan bool)

	localBarber := &BarberReader{
		Stop:      stop,
		Start:     start,
		TimeSlept: timeSlept,
		IsBusy:    isBusy,
		ID:        id,
		Log:       log,
		End:       end,
		Kill:      kill,
	}

	barber := &Barber{
		id:        id,
		log:       new(SWriter),
		busy:      false,
		Stop:      stop,
		Start:     start,
		TimeSlept: timeSlept,
		IsBusy:    isBusy,
		Log:       log,
		End:       end,
		Kill:      kill,
	}

	// Start barber go routine
	go barber.GoLive()

	return localBarber
}

// BestBarber finds the barber who has been waiting the longest
func BestBarber(barbers []*BarberReader) (*BarberReader, error) {
	bestTime := time.Duration(0)
	var bestBarber *BarberReader = nil

	t := make(chan time.Duration)
	b := make(chan bool)
	for _, barber := range barbers {
		barber.IsBusy <- b
		if !<-b {
			barber.TimeSlept <- t
			newTime := <-t
			if newTime > bestTime {
				bestTime = newTime
				bestBarber = barber
			}
		}
	}

	if len(barbers) == 0 {
		return nil, fmt.Errorf("No barbers in list")
	}

	if bestBarber != nil {
		return bestBarber, nil
	}

	return nil, fmt.Errorf("All barbers are busy")
}

// allBarbersFinished checks if all of the barbers are free (not busy)
func allBarbersFinished(barbers []*BarberReader) bool {
	for _, barber := range barbers {
		c := make(chan bool)
		barber.IsBusy <- c
		busy := <-c
		if busy {
			return false
		}
	}

	return true
}
