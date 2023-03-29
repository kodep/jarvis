package filters

import (
	"context"
	"reflect"
	"regexp"

	"github.com/kodep/jarvis/internal/mattermost/events"
	"github.com/kodep/jarvis/internal/mattermost/handlers"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/samber/lo"
)

func ByEvent(v events.Event) handlers.FilterFn[events.Event] {
	return func(ctx context.Context, e events.Event) (bool, error) {
		return reflect.TypeOf(e) == reflect.TypeOf(v), nil
	}
}

func ByEventType(eventType events.EventType) handlers.FilterFn[events.Event] {
	return func(ctx context.Context, e events.Event) (bool, error) {
		return e.EventType() == eventType, nil
	}
}

func ByRegexp(r *regexp.Regexp) handlers.FilterFn[*events.PostEvent] {
	return func(ctx context.Context, e *events.PostEvent) (bool, error) {
		return r.MatchString(e.Message), nil
	}
}

func ByChannelID(channelID string) handlers.FilterFn[*events.PostEvent] {
	return func(ctx context.Context, e *events.PostEvent) (bool, error) {
		return e.ChannelId == channelID, nil
	}
}

func MentionedMe(me *model.User) handlers.FilterFn[*events.PostEvent] {
	return func(ctx context.Context, e *events.PostEvent) (bool, error) {
		return lo.Some(e.Mentions, []string{me.Id}), nil
	}
}
