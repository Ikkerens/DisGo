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
	if eventType.ConvertibleTo(eventInterface) {
		panic("The second argument of the passed Func should be a known Event.")
	}

	// Create a zero'd instance of the particular Event, so that we can call eventName() on it
	eventInstance := reflect.New(eventType).Interface()
	zeroEvent := eventInstance.(Event)

	// Wrap the event handler in a reflected function that "type asserts" the event
	wrapper := func(session *Session, event *Event) {
		go handler.Call([]reflect.Value{reflect.ValueOf(session), reflect.ValueOf(*event).Elem().Convert(eventType)})
	}

	list, exists := handlers[zeroEvent.eventName()]
	if !exists {
		handlers[zeroEvent.eventName()] = []eventHandler{wrapper}
	} else {
		handlers[zeroEvent.eventName()] = append(list, wrapper)
	}
}

func (s *Session) dispatchEvent(frame *receivedFrame) {
	var event Event

	switch frame.EventName {
	case "READY":
		event = &ReadyEvent{}
	case "RESUMED":
		event = &ResumedEvent{}
	case "GUILD_CREATE":
		event = &GuildCreateEvent{}
	case "MESSAGE_CREATE":
		event = &MessageCreateEvent{}
	case "MESSAGE_DELETE":
		event = &MessageDeleteEvent{}
	case "TYPING_START":
		event = &TypingStartEvent{}
	default:
		logger.Errorf("Event with name '%s' was dispatched by Discord, but we don't know this event. (DisGo outdated?)", frame.EventName)
		return
	}

	err := json.Unmarshal(frame.Data, &event)
	if err != nil {
		logger.ErrorE(err)
		return
	}

	event.setSession(s)

	handlerSlice, exists := handlers[event.eventName()]
	if exists {
		for _, handler := range handlerSlice {
			handler(s, &event)
		}
	}
}
