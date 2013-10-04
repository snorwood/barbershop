package barbershop

import (
	"testing"
	"time"
)

func TestQueueEnqueue(t *testing.T) {
	queue := NewGeneralQueue(10)

	for i := 0; i < 10; i++ {
		err := queue.Enqueue(i)
		if err != nil {
			t.Error("Unexpected Error: ", err)
		}
	}

	if queue.Enqueue(11) == nil {
		t.Error("Ignored queue capacity")
	}
}

func TestQueueDequeue(t *testing.T) {
	queue := NewGeneralQueue(10)

	for i := 0; i < 10; i++ {
		err := queue.Enqueue(i)
		if err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < 10; i++ {
		value, err := queue.Dequeue()
		if value != i {
			t.Error("Dequeued value is not correct", value)
		}

		if err != nil {
			t.Error("Unexpected Error: ", err)
		}
	}

	_, err := queue.Dequeue()
	if err == nil {
		t.Error("Dequeued from an empty list without error")
	}
}

func TestQueuePeek(t *testing.T) {
	queue := NewGeneralQueue(10)

	for i := 0; i < 10; i++ {
		err := queue.Enqueue(i)
		if err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < 10; i++ {
		value, err := queue.Dequeue()
		if value != i {
			t.Error("Dequeued value is not correct", value)
		}

		if err != nil {
			t.Error("Unexpected Error: ", err)
		}
	}

	_, err := queue.Dequeue()
	if err == nil {
		t.Error("Dequeued from an empty list without error")
	}
}

func TestSendSubscriber(t *testing.T) {
	subscriber := NewSubscriber()
	ch := make(chan Subscriber)
	go SendSubscriber(subscriber, ch)

	select {
	case reverseSubscriber := <-ch:
		go func() {
			reverseSubscriber.Send <- "test"
			reverseSubscriber.StopReceiving <- true
			<-reverseSubscriber.Receive
			<-reverseSubscriber.StopSending
		}()
	case <-time.After(time.Second):
		t.Error("Subscriber was not sent or took too long")
	}

	select {
	case <-subscriber.Receive:
	case <-time.After(time.Second):
		t.Error("Message was not sent on Send channel or took too long")
	}

	select {
	case <-subscriber.StopSending:
	case <-time.After(time.Second):
		t.Error("Message was not sent on StopReceiving channel or took too long")
	}

	select {
	case subscriber.Send <- "test":
	case <-time.After(time.Second):
		t.Error("Message was not recieved on Receive channel or took too long")
	}

	select {
	case subscriber.StopReceiving <- true:
	case <-time.After(time.Second):
		t.Error("Message was not received on StopSending channel or took too long")
	}
}
