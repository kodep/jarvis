package events

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v6/model"
)

type (
	UnknownEventError struct {
		WsEvent *model.WebSocketEvent
	}

	DecodingError struct {
		WsEvent *model.WebSocketEvent
		Cause   error
	}
)

type eventDecoder interface {
	Accept(*model.WebSocketEvent) bool
	Decode(*model.WebSocketEvent) (Event, error)
}

var _decoders = [2]eventDecoder{ //nolint:gochecknoglobals // it's ok
	postEventDecoder{},
	unknownEventDecoder{}, // should be always last
}

func decodeEvent(e *model.WebSocketEvent) (Event, error) {
	for _, decoder := range _decoders {
		if accepted := decoder.Accept(e); !accepted {
			continue
		}

		d, err := decoder.Decode(e)
		if err != nil {
			return nil, &DecodingError{WsEvent: e, Cause: err}
		}

		return d, nil
	}

	return nil, &UnknownEventError{WsEvent: e}
}

func (e *UnknownEventError) Error() string {
	return fmt.Sprintf("unknown event type: %s", e.WsEvent.EventType())
}

func (e *DecodingError) Error() string {
	return "decoding error: " + e.Cause.Error()
}
