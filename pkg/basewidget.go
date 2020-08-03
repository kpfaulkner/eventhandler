package pkg

import (
	log "github.com/sirupsen/logrus"
)

// BaseWidget is the common element of ALL disable items. Widgets, buttons, panels etc.
// Should only handle some basic items like location, width, height and possibly events.
type EventHandlerManager struct {
	// These are other widgets/components that are listening to THiS widget. Ie we will broadcast to them!
	eventListeners map[int][]chan IEvent

	// incoming events to THIS widget (ie stuff we're listening to!)
	incomingEvents chan IEvent
}

func NewEventHandlerManager() *EventHandlerManager {

	eh := EventHandlerManager{}
	eh.eventListeners = make(map[int][]chan IEvent)
	eh.incomingEvents = make(chan IEvent, 1000) // too much?

	return &eh
}

func (e *EventHandlerManager) AddEventListener(eventType int, ch chan IEvent) error {
	if _, ok := e.eventListeners[eventType]; ok {
		e.eventListeners[eventType] = append(e.eventListeners[eventType], ch)
	} else {
		e.eventListeners[eventType] = []chan IEvent{ch}
	}

	return nil
}

func (e *EventHandlerManager) RemoveEventListener(eventType int, ch chan IEvent) error {
	if _, ok := e.eventListeners[eventType]; ok {
		for i := range e.eventListeners[eventType] {
			if e.eventListeners[eventType][i] == ch {
				e.eventListeners[eventType] = append(e.eventListeners[eventType][:i], e.eventListeners[eventType][i+1:]...)
				break
			}
		}
	}
	return nil
}

// Emit event for  all listeners to receive
func (e *EventHandlerManager) EmitEvent(event IEvent) error {

	eventToUse := event

	if _, ok := e.eventListeners[eventToUse.EventType()]; ok {
		for _, handler := range e.eventListeners[eventToUse.EventType()] {
			go func(handler chan IEvent) {
				handler <- eventToUse
			}(handler)
		}
	}

	return nil
}

func (e *EventHandlerManager) ListenToIncomingEvents() error {

	for {
		ev := <-e.incomingEvents

		log.Debugf("event %v", ev)
		/*
			// do our local event processing (HandleEvent) then pass onto other listeners (assuming order would be important here).
			used, err := e.eventHandler(ev)
			if err != nil {
				log.Errorf("Unable to HandleEvent from widget: %s", err.Error())
				continue
			}

			// if USED by this widget... then pass it onto the child widgets.
			// if NOT used by this widget.... its nothing to do with us... dont
			// propagate.
			if used {
				// if mouse event, convert to local co-ord system?
				err := b.EmitEvent(ev)
				if err != nil {
					log.Errorf("Unable to emit event from widget: %s", err.Error())
					// wont break out here... assuming/hoping that this is just a once off :)
				}
			} */
	}
	return nil
}
