package pkg

const (
	EventType1 int = iota
	EventType2
)

type IEvent interface {
	Name() string
	EventType() int
	Action() error
}

type Event struct {
	eventName string
	eventType int
}

func NewEvent(eventType int) Event {
	e := Event{}
	return e
}
