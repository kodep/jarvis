package filters

import (
	"regexp"

	"github.com/kodep/jarvis/internal/mattermost/events"
	"github.com/kodep/jarvis/internal/mattermost/handlers"
	"github.com/samber/lo"
)

func ByEventType(eventType events.EventType) handlers.FilterFn[events.Event] {
	return func(ctx handlers.Context, e events.Event) (bool, error) {
		return e.EventType() == eventType, nil
	}
}

func ByRegexp(r *regexp.Regexp) handlers.FilterFn[events.PostEvent] {
	return func(ctx handlers.Context, e events.PostEvent) (bool, error) {
		return r.MatchString(e.Message()), nil
	}
}

func ByChannelID(channelID string) handlers.FilterFn[events.PostEvent] {
	return func(ctx handlers.Context, e events.PostEvent) (bool, error) {
		return e.ChannelID() == channelID, nil
	}
}

func MentionedMe() handlers.FilterFn[events.PostEvent] {
	return func(ctx handlers.Context, e events.PostEvent) (bool, error) {
		me := ctx.Client().User()
		return lo.Some(e.Mentions(), []string{me.Id}), nil
	}
}

func Conjunction[T events.Event](fns ...handlers.FilterFn[T]) handlers.FilterFn[T] {
	return func(ctx handlers.Context, e T) (bool, error) {
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

func Disjunction[T events.Event](fns ...handlers.FilterFn[T]) handlers.FilterFn[T] {
	return func(ctx handlers.Context, e T) (bool, error) {
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
