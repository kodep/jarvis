package filters

import (
	"context"
	"regexp"

	"github.com/kodep/jarvis/internal/mattermost/events"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/samber/lo"
)

type FilterFn[T events.Event] func(ctx context.Context, e T) (bool, error)

func ByEventType(eventType events.EventType) FilterFn[events.Event] {
	return func(ctx context.Context, e events.Event) (bool, error) {
		return e.EventType() == eventType, nil
	}
}

func ByRegexp(r *regexp.Regexp) FilterFn[events.PostEvent] {
	return func(ctx context.Context, e events.PostEvent) (bool, error) {
		return r.MatchString(e.Message()), nil
	}
}

func ByChannelID(channelID string) FilterFn[events.PostEvent] {
	return func(ctx context.Context, e events.PostEvent) (bool, error) {
		return e.ChannelID() == channelID, nil
	}
}

func MentionedMe(me *model.User) FilterFn[events.PostEvent] {
	return func(ctx context.Context, e events.PostEvent) (bool, error) {
		return lo.Some(e.Mentions(), []string{me.Id}), nil
	}
}

func Conjunction[T events.Event](fns ...FilterFn[T]) FilterFn[T] {
	return func(ctx context.Context, e T) (bool, error) {
		var (
			ok  = true
			err error
		)

		for _, fn := range fns {
			ok, err = fn(ctx, e)
			if err != nil {
				return false, err
			}

			if !ok {
				break
			}
		}

		return ok, nil
	}
}

func Disjunction[T events.Event](fns ...FilterFn[T]) FilterFn[T] {
	return func(ctx context.Context, e T) (bool, error) {
		var (
			ok  = true
			err error
		)

		for _, fn := range fns {
			ok, err = fn(ctx, e)
			if err != nil {
				return false, err
			}

			if ok {
				break
			}
		}

		return ok, nil
	}
}
