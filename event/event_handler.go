package event

import "fmt"

var queue chan EventState = make(chan EventState)

func Queue(es EventState) {
	queue <- es
}

func Shutdown() {
	close(queue)
}

func Listen() {
	for es := range queue {
		handleEvent(&es)
	}
}

func handleEvent(es *EventState) {
	et := eventType(es)
	switch et {
	// TODO, handle lots of event types ...
	default:
		fmt.Printf("e: %s %s\n", et, es.String())
	}
}

func eventType(es *EventState) EventType {
	// TODO, handle lots of event states ...
	for chord, et := range standard {
		if es.matches(chord) {
			return et
		}
	}
	return Evt_
}
