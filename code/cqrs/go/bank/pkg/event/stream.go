package event

type Stream interface {
	Load() <-chan Event
}
