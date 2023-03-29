package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/wire"
	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"github.com/kodep/jarvis/internal/mattermost/events"
	"github.com/kodep/jarvis/internal/mattermost/filters"
	"github.com/kodep/jarvis/internal/mattermost/guards"
	"github.com/kodep/jarvis/internal/mattermost/handlers"
	"github.com/kodep/jarvis/internal/oboobs"
	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

type EventsHandler func(ctx context.Context, e events.Event) error

type (
	LogsHandler  handlers.MiddlwareFn[events.Event]
	BoobsHandler handlers.MiddlwareFn[events.Event]
	ButtsHandler handlers.MiddlwareFn[events.Event]
)

//nolint:gochecknoglobals // it's needed for wire
var EventsHandlersSet = wire.NewSet(
	ProvideEventsHandler,
	ProvideLogsHandler,
	ProvideBoobsHandler,
	ProvideButtsHandler,
)

func ProvideEventsHandler(logsHandler LogsHandler, boobsHandler BoobsHandler, buttsHandler ButtsHandler) EventsHandler {
	return handlers.TerminatePipe(handlers.Pipe(
		(handlers.MiddlwareFn[events.Event])(logsHandler),
		(handlers.MiddlwareFn[events.Event])(boobsHandler),
		(handlers.MiddlwareFn[events.Event])(buttsHandler),
	))
}

func ProvideLogsHandler(logger *zap.Logger) LogsHandler {
	return LogsHandler(func(ctx context.Context, e events.Event, next handlers.NextFn[events.Event]) error {
		logger.Debug("Received event", zap.String("EventType", e.RawEventType()))

		err := next(handlers.WithLogger(ctx, logger), e)

		if !e.Acknowledged() && err == nil {
			logger.Debug("Skipped event", zap.String("EventType", e.RawEventType()))
		}

		return err
	})
}

//nolint:dupl // todo
func ProvideBoobsHandler(client *mattermost.Client, boobs oboobs.BoobsClient, conf Config) BoobsHandler {
	regex := regexp.MustCompile("show boobs")

	return BoobsHandler(guards.EventGuard(
		handlers.Filter(filters.ByChannelID(conf.BoobsChannelID), handlers.Filter(filters.ByRegexp(regex),
			func(ctx context.Context, e *events.PostEvent, next handlers.NextFn[*events.PostEvent]) error {
				logger := handlers.GetLogger(ctx)

				logger.Debug("Asked me to show boobs")

				boob, err := boobs.Random(ctx)
				if err != nil {
					return fmt.Errorf("failed to get boobs: %w", err)
				}

				logger.Debug("Send boobs", zap.String("boobs", boob.URL))

				post := &model.Post{ChannelId: e.ChannelId, Message: boob.URL}

				if _, err = client.SendPost(post); err != nil {
					return fmt.Errorf("failed to answer: %w", err)
				}

				e.Ack()
				return nil
			},
		))))
}

//nolint:dupl // todo
func ProvideButtsHandler(client *mattermost.Client, butts oboobs.ButtsClient, conf Config) ButtsHandler {
	regex := regexp.MustCompile("show butts")

	return ButtsHandler(guards.EventGuard(
		handlers.Filter(filters.ByChannelID(conf.BoobsChannelID), handlers.Filter(filters.ByRegexp(regex),
			func(ctx context.Context, e *events.PostEvent, next handlers.NextFn[*events.PostEvent]) error {
				logger := handlers.GetLogger(ctx)

				logger.Debug("Asked me to show butts")

				butt, err := butts.Random(ctx)
				if err != nil {
					return fmt.Errorf("failed to get butts: %w", err)
				}

				logger.Debug("Send butts", zap.String("butts", butt.URL))

				post := &model.Post{ChannelId: e.ChannelId, Message: butt.URL}

				if _, err = client.SendPost(post); err != nil {
					return fmt.Errorf("failed to answer: %w", err)
				}

				e.Ack()
				return nil
			},
		))))
}
