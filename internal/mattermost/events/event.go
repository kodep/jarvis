package events

import (
	"github.com/mattermost/mattermost/server/public/model"
)

type EventType string

type Event interface {
	Acknowledged() bool
	Ack()
	NoAck()
	EventType() EventType
	RawEventType() string
	RawData() map[string]interface{}
}

func NewEvent(e *model.WebSocketEvent) (Event, error) {
	return decodeEvent(e)
}
