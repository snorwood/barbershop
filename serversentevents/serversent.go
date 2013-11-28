package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type client struct {
	Recieve chan string
	ID      int
}

type Broker struct {
	clients        map[int]client
	newClients     chan client
	defunctClients chan client
	messages       chan string
	maxID          int
}

func (self *Broker) Start() {
	servicing := true
	self.maxID = 0
	go func() {
		for servicing {
			select {
			case newClient := <-self.newClients:
				self.clients[newClient.ID] = newClient
				log.Println("Added new client")

			case defunctClient := <-self.defunctClients:
				delete(self.clients, defunctClient.ID)
				log.Println("Removed Client")

			case msg := <-self.messages:
				for _, client := range self.clients {
					client.Recieve <- msg
				}

				log.Printf("Broadcast \"%s\" to %d clients", msg, len(self.clients))
			}
		}
	}()
}

func (self *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	cn, ok := w.(http.CloseNotifier)
	if !ok {
		http.Error(w, "Close notifier unsupported", http.StatusInternalServerError)
	}
	closeConnection := cn.CloseNotify()

	joiningClient := client{ID: self.maxID, Recieve: make(chan string)}
	self.maxID += 1

	self.newClients <- joiningClient
	defer func() {
		self.defunctClients <- joiningClient
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	connectionOpen := true

	for connectionOpen {
		select {
		case msg := <-joiningClient.Recieve:
			fmt.Fprintf(w, "data: Message: %s\n\n", msg)
			f.Flush()
		case <-closeConnection:
			connectionOpen = false
		}
	}

	log.Println("Finished HTTP request at", r.URL.Path)
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Error parsing template")
	}

	t.Execute(w, "Steven")

	log.Println("Finished HTTP request at ", r.URL.Path)
}

func main() {
	b := &Broker{
		make(map[int]client),
		make(chan client),
		make(chan client),
		make(chan string),
		0,
	}

	b.Start()
	http.Handle("/events/", b)
	go func() {
		for i := 0; ; i++ {
			b.messages <- fmt.Sprintf("%d - the time is %v", i, time.Now().Format("Jan 2 15:04"))
			log.Printf("Sent message %d ", i)
			time.Sleep(5 * time.Second)
		}
	}()

	http.Handle("/", http.HandlerFunc(MainPageHandler))
	http.ListenAndServe(":8000", nil)
}
