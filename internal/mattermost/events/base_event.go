package events

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

type BaseEvent struct {
	eventType EventType

	wsEvent      *model.WebSocketEvent
	acknowledged bool
}

var _ Event = (*BaseEvent)(nil)

func newBaseEvent(e *model.WebSocketEvent, eventType EventType) BaseEvent {
	return BaseEvent{
		eventType: eventType,
		wsEvent:   e,
	}
}

func (e *BaseEvent) Acknowledged() bool {
	return e.acknowledged
}

func (e *BaseEvent) Ack() {
	e.acknowledged = true
}

func (e *BaseEvent) NoAck() {
	e.acknowledged = false
}

func (e *BaseEvent) EventType() EventType {
	return e.eventType
}

func (e *BaseEvent) RawEventType() string {
	return e.wsEvent.EventType()
}

func (e *BaseEvent) RawData() map[string]interface{} {
	return e.wsEvent.GetData()
}
