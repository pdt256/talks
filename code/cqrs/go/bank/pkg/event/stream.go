package event

import "io"

type Stream interface {
	Load(io.Reader) <-chan Event
	Save(io.Writer, <-chan Event) <-chan error
}
