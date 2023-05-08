package subscribe

import "test/models"

type Subsribe struct {
	Ready  chan chan models.Event
	Cancel chan chan models.Event
}

func (this *Subsribe) Start() (eventChan <-chan models.Event) {
	eventChan = make(chan models.Event, 100)

	this.Ready <- eventChan
	return
}

func (this *Subsribe) End(rottenEventChan chan models.Event) {
	this.Cancel <- rottenEventChan
}
