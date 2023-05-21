package saver

import (
	"sync"
	"trading/common/connectors"
	"trading/storage"
)

type Saver struct {
	PostmanPort      uint
	Config           SaverConfig
	PostmanConnector connectors.PostmanConnector
}

func MakeSaver(config SaverConfig, postmanPort uint) *Saver {
	return &Saver{
		PostmanPort:      postmanPort,
		Config:           config,
		PostmanConnector: connectors.MakePostmanConnector(postmanPort)}
}

func (this *Saver) Listen(instrument string) {
	events_keeper := storage.MakeEventsKeeper(instrument)
	for event := range this.PostmanConnector.SubscribeDepth(instrument) {
		events_keeper.Save(event)
	}
}

func (this *Saver) Run() {

	wg := sync.WaitGroup{}

	for _, instrument := range this.Config.SubscribedInstruments {

		// Узнаём порт по которому нужно подключаться
		this.PostmanConnector.AddInstrument(instrument)

		// запускаем подписку по этому инструменту
		go func(instrument string) {
			wg.Add(1)
			this.Listen(instrument)
			wg.Done()
		}(instrument)
	}

	wg.Wait()
}
