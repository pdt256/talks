package event

import (
	"reflect"
)

type Serializer interface {
	Bind(events ...Event)
	Serialize(event Event) ([]byte, error)
	Deserialize(serializedData []byte, eventTypeName string) (Event, error)
}

func Type(event Event) (string, reflect.Type) {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name(), t
}
