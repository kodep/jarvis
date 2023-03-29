package events

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

const (
	UnknownEventType EventType = "unknown"
)

type UnknownEvent struct {
	BaseEvent
}

type unknownEventDecoder struct{}

var _ eventDecoder = (*unknownEventDecoder)(nil)

func (d unknownEventDecoder) Accept(e *model.WebSocketEvent) bool {
	// unknown decoder accepts all events
	return true
}

func (d unknownEventDecoder) Decode(e *model.WebSocketEvent) (Event, error) {
	return &UnknownEvent{
		BaseEvent: newBaseEvent(e, UnknownEventType),
	}, nil
}
