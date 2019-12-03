2019-10-24 - Event Modeling - From Sticky Notes to Go

20 min presentation for the [West LA Go Meetup](https://www.meetup.com/West-LA-Go/events/265413849/).

This talk covered three topics: Event Storming, a CQRS architecture with Event Sourcing, and implementing Aggregates and Projections in Go following TDD.

We started by Event Storming a bank account behavior workflow with sticky notes generating Events, Commands, and Constraints. We then turn these workflows into unit tests matching one-to-one with the expected domain behavior:

* **Given** these events
* **When** this command is executed
* **Then** we expect these events to be emitted.

Finally we dove deep into example code, in the Go language, to see the interaction between the event store, aggregates, projections, and business logic constraints.

 
* Slides: [event-modeling-from-sticky-notes-to-go.pdf](event-modeling-from-sticky-notes-to-go.pdf)
* Video: https://youtu.be/i7_edqzneyM
* Code: [https://github.com/pdt256/talks/tree/master/code/cqrs/go/bank](https://github.com/pdt256/talks/tree/dd1e8c2b903fb8b7ca55115465b47458cf49b1ac/code/cqrs/go/bank)


[![Event Modeling - From Sticky Notes to Go](https://github.com/pdt256/talks/raw/master/2019-10-24-event-modeling-from-sticky-notes-to-go/photos/screenshot.jpg)](https://youtu.be/i7_edqzneyM)

[![Event Modeling - From Sticky Notes to Go - Event Storming](https://github.com/pdt256/talks/raw/master/2019-10-24-event-modeling-from-sticky-notes-to-go/photos/slide-event-storming.jpg)](event-modeling-from-sticky-notes-to-go.pdf)
[![Event Modeling - From Sticky Notes to Go - CQRS + Event Souring](https://github.com/pdt256/talks/raw/master/2019-10-24-event-modeling-from-sticky-notes-to-go/photos/slide-cqrs-es.jpg)](event-modeling-from-sticky-notes-to-go.pdf)
