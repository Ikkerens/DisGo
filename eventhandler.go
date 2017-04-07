package disgo

import (
	"encoding/json"
	"reflect"

	"github.com/slf4go/logger"
)

type eventHandler func(session *Session, event Event)

var (
	eventInterface reflect.Type
	handlers       map[string][]eventHandler
)

func init() {
	eventInterface = reflect.TypeOf((*Event)(nil)).Elem()
	handlers = make(map[string][]eventHandler)
}

func (s *Session) RegisterEventHandler(handlerI interface{}) {
	// These panics should be purely informational
	defer logger.RecoverStack()

	handler := reflect.ValueOf(handlerI)
	handlerType := handler.Type()

	// Is the passed handler a Func?
	if handler.Kind() != reflect.Func {
		panic("Passed zeroEvent handler is not a Func.")
	}
	// The signature requires two arguments, do we match that
	if handlerType.NumIn() != 2 {
		panic("Passed zeroEvent handler should be a Func with 2 arguments")
	}
	// Is the first argument a Session pointer?
	if handlerType.In(0) != reflect.TypeOf(s) {
		panic("The first argument of the passed Func should be of Type *disgo.Session")
	}

	// Is the second argument a struct that implements Event so we can obtain the Event Name?
	eventType := handlerType.In(1)
	if !eventType.Implements(eventInterface) {
		panic("The second argument of the passed Func should be a known zeroEvent pointer.")
	}

	// Create a zero'd instance of the particular Event, so that we can call EventName() on it
	eventInstance := reflect.New(eventType).Interface()
	zeroEvent := eventInstance.(Event)

	wrapper := func(session *Session, event Event) {
		handler.Call([]reflect.Value{reflect.ValueOf(session), reflect.ValueOf(event).Elem().Convert(eventType)})
	}

	list, exists := handlers[zeroEvent.EventName()]
	if !exists {
		handlers[zeroEvent.EventName()] = []eventHandler{wrapper}
	} else {
		handlers[zeroEvent.EventName()] = append(list, wrapper)
	}
}

func (s *Session) dispatchEvent(frame *receivedFrame) {
	var event Event

	switch frame.EventName {
	case "READY":
		event = &ReadyEvent{}
	default:
		logger.Errorf("Event with name '%s' was dispatched by Discord, but we don't know this event. (DisGo outdated?)", frame.EventName)
		return
	}

	err := json.Unmarshal(frame.Data, &event)
	if err != nil {
		logger.ErrorE(err)
		return
	}

	handlerSlice, exists := handlers[event.EventName()]
	if exists {
		for _, handler := range handlerSlice {
			handler(s, event)
		}
	}
}
