package disgo

import (
	"encoding/json"
	"reflect"

	"github.com/slf4go/logger"
)

type Event interface {
	eventName() string

	setSession(*Session)
}

type eventHandler func(session *Session, event *Event)

var (
	eventInterface reflect.Type
	handlers       map[string][]eventHandler
)

func init() {
	eventInterface = reflect.TypeOf((*Event)(nil))
	handlers = make(map[string][]eventHandler)
}

func (s *Session) RegisterEventHandler(handler interface{}) {
	s.registerEventHandler(handler, true)
}

func (s *Session) registerEventHandler(handlerI interface{}, goroutine bool) {
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
	if eventType.ConvertibleTo(eventInterface) {
		panic("The second argument of the passed Func should be a known Event.")
	}

	// Wrap the event handler in a reflected function that "type asserts" the event
	wrapper := func(session *Session, event *Event) {
		handler.Call([]reflect.Value{reflect.ValueOf(session), reflect.ValueOf(*event).Elem().Convert(eventType)})
	}

	// Wrap it in a goroutine if needed
	if goroutine {
		internal := wrapper
		wrapper = func(session *Session, event *Event) {
			go internal(session, event)
		}
	}

	// Create a zero'd instance of the particular Event, so that we can call eventName() on it
	eventName := reflect.New(eventType).Interface().(Event).eventName()
	list, exists := handlers[eventName]
	if !exists {
		handlers[eventName] = []eventHandler{wrapper}
	} else {
		handlers[eventName] = append(list, wrapper)
	}
}

func (s *Session) dispatchEvent(frame *receivedFrame) {
	event := allocateEvent(frame.EventName)

	if event == nil {
		return
	}

	if err := json.Unmarshal(frame.Data, &event); err != nil {
		logger.ErrorE(err)
		return
	}

	logger.Debugf("Dispatching event %s to handlers", (*event).eventName())
	(*event).setSession(s)
	if handlerSlice, exists := handlers[(*event).eventName()]; exists {
		for _, handler := range handlerSlice {
			handler(s, event)
		}
	}
}
