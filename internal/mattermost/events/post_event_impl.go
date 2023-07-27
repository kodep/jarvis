package events

import (
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

type postEventImpl struct {
	BaseEvent
	post        model.Post
	channelName string
	mentions    []string
}

var _ PostEvent = (*postEventImpl)(nil)

func (e *postEventImpl) ID() string {
	return e.post.Id
}

func (e *postEventImpl) CreatedAt() time.Time {
	return time.Unix(e.post.CreateAt, 0)
}

func (e *postEventImpl) UpdatedAt() time.Time {
	return time.Unix(e.post.UpdateAt, 0)
}

func (e *postEventImpl) EditedAt() time.Time {
	return time.Unix(e.post.EditAt, 0)
}

func (e *postEventImpl) UserID() string {
	return e.post.ChannelId
}

func (e *postEventImpl) ChannelID() string {
	return e.post.ChannelId
}

func (e *postEventImpl) Message() string {
	return e.post.Message
}

func (e *postEventImpl) Mentions() []string {
	return e.mentions
}
