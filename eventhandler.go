package disgo

import (
	"encoding/json"
	"reflect"

	"github.com/slf4go/logger"
)

type eventHandler func(*Session, *Event)

var (
	eventInterface reflect.Type
	handlers       map[string][]eventHandler
)

func init() {
	eventInterface = reflect.TypeOf((*Event)(nil)).Elem()
	handlers = make(map[string][]eventHandler)
}

func (s *Session) RegisterEventHandler(handler interface{}) {
	// These panics should be purely informational
	defer logger.RecoverStack()

	// Is the passed handler a Func?
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("Passed event handler is not a Func.")
	}
	// The signature requires two arguments, do we match that
	if handlerType.NumIn() != 2 {
		panic("Passed event handler should be a Func with 2 arguments")
	}
	// Is the first argument a Session pointer?
	if handlerType.In(0) != reflect.TypeOf(s) {
		panic("The first argument of the passed Func should be of Type *disgo.Session")
	}

	// Is the second argument a struct that implements Event so we can obtain the Event Name?
	eventType := handlerType.In(1)
	if !eventType.Implements(eventInterface) {
		panic("The second argument of the passed Func should be a known event pointer.")
	}

	// Create a zero'd instance of the particular Event, so that we can call EventName() on it
	eventInstance := reflect.New(eventType.Elem()).Interface()
	event := eventInstance.(Event)

	logger.Infof("Registered event: %s", event.EventName())
	list, exists := handlers[event.EventName()]
	if !exists {
		handlers[event.EventName()] = []eventHandler{handler.(eventHandler)}
	} else {
		handlers[event.EventName()] = append(list, handler.(eventHandler))
	}
}

func dispatchEvent(frame *receivedFrame) {
	var event interface{}

	switch frame.EventName {

	case "READY":
		event = ReadyEvent{}

	}

	json.Unmarshal(frame.Data, &event)
}
