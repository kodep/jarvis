package events

import (
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

const (
	PostedEventType EventType = EventType(model.WebsocketEventPosted)
)

type PostEvent interface {
	Event
	ID() string
	CreatedAt() time.Time
	UpdatedAt() time.Time
	EditedAt() time.Time
	UserID() string
	ChannelID() string
	Message() string
	Mentions() []string
}
