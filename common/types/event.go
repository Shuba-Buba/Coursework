package types

import "time"

type EventType int

const (
	EventTypeOrderbookUpdate EventType = iota
	EventTypeTrades          EventType = iota
	EventTypeSnapshot        EventType = iota
)

type Event struct {
	Timestamp time.Time
	Type      EventType
	Data      string
}
