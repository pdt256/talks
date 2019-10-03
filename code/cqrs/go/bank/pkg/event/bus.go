package event

type Bus interface {
	Subscribe(subscriber Subscriber)
	Publish(data ...Event)
}

type Subscriber interface {
	Accept(event Event)
}
