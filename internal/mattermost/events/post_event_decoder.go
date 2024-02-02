package events

import (
	"encoding/json"
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
)

type postEventDecoder struct{}

var _ eventDecoder = (*postEventDecoder)(nil)

func (d postEventDecoder) Accept(e *model.WebSocketEvent) bool {
	return e.EventType() == model.WebsocketEventPosted
}

func (d postEventDecoder) Decode(e *model.WebSocketEvent) (Event, error) {
	const (
		channelNameKey = "channel_name"
		mentionsKey    = "mentions"
		postKey        = "post"
	)

	data := e.GetData()
	post := &postEventImpl{
		BaseEvent: newBaseEvent(e, UnknownEventType),
	}

	if str, hasChannelNameKey := data[channelNameKey].(string); hasChannelNameKey {
		post.channelName = str
	}

	if str, hasMentionsKey := data[mentionsKey].(string); hasMentionsKey {
		if err := json.Unmarshal([]byte(str), &post.mentions); err != nil {
			return nil, fmt.Errorf("failed to parse event: %w", err)
		}
	}

	if str, hasPostKey := data[postKey].(string); hasPostKey {
		if err := json.Unmarshal([]byte(str), &post.post); err != nil {
			return nil, fmt.Errorf("failed to parse event: %w", err)
		}
	} else {
		return nil, fmt.Errorf("failed to parse event: missing post")
	}

	return post, nil
}
