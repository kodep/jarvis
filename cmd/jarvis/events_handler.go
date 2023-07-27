package main

import (
	"fmt"
	"regexp"

	"github.com/google/wire"
	"github.com/kodep/jarvis/internal/mattermost/events"
	"github.com/kodep/jarvis/internal/mattermost/factories"
	"github.com/kodep/jarvis/internal/mattermost/filters"
	"github.com/kodep/jarvis/internal/mattermost/guards"
	"github.com/kodep/jarvis/internal/mattermost/handlers"
	"github.com/kodep/jarvis/internal/oboobs"
	"github.com/mattermost/mattermost/server/public/model"
	"go.uber.org/zap"
)

type EventsHandler func(ctx handlers.Context, e events.Event) error

type (
	logsHandler handlers.MiddlwareFn[events.Event]
)

//nolint:gochecknoglobals // it's needed for wire
var EventsHandlersSet = wire.NewSet(
	provideEventsHandler,
	provideLogsHandler,
	wire.Struct(new(boobsAndButtsHandler), "*"),
)

func provideEventsHandler(
	conf Config,
	logsHandler logsHandler,
	bb boobsAndButtsHandler,
) EventsHandler {
	boobsRegexp := regexp.MustCompile("(?i)show boobs")
	buttsRegexp := regexp.MustCompile("(?i)show butts")

	boobsAndButts := handlers.Filter(filters.ByChannelID(conf.BoobsChannelID), handlers.Pipe(
		handlers.Filter(filters.ByRegexp(boobsRegexp), bb.handleBoobs),
		handlers.Filter(filters.ByRegexp(buttsRegexp), bb.handleButts),
	))

	return handlers.TerminatePipe(handlers.Pipe(
		(handlers.MiddlwareFn[events.Event])(logsHandler),
		guards.EventGuard(boobsAndButts),
	))
}

func provideLogsHandler() logsHandler {
	return func(ctx handlers.Context, e events.Event, next handlers.NextFn[events.Event]) error {
		ctx.Logger().Debug("Received event", zap.String("EventType", e.RawEventType()))

		err := next(ctx, e)

		if !e.Acknowledged() && err == nil {
			ctx.Logger().Debug("Skipped event", zap.String("EventType", e.RawEventType()))
		}

		return err
	}
}

type boobsAndButtsHandler struct {
	Boobs oboobs.BoobsClient
	Butts oboobs.ButtsClient
}

func (bh boobsAndButtsHandler) handleBoobs( //nolint:dupl // 😕
	ctx handlers.Context,
	e events.PostEvent,
	next handlers.NextFn[events.PostEvent],
) error {
	logger := ctx.Logger().With(zap.String("handler", "boobs"))

	logger.Debug("Asked me to show boobs")

	factories.SendTyping(ctx.WSClient(), e.ChannelID())

	boob, err := bh.Boobs.Random(ctx.Context())
	if err != nil {
		return fmt.Errorf("failed to get boobs: %w", err)
	}

	logger.Debug("Send boobs", zap.String("boobs", boob.URL))

	post := &model.Post{ChannelId: e.ChannelID(), Message: boob.URL}

	if _, err = ctx.Client().SendPost(ctx.Context(), post); err != nil {
		return fmt.Errorf("failed to answer: %w", err)
	}

	e.Ack()
	return nil
}

func (bh boobsAndButtsHandler) handleButts( //nolint:dupl // 😕
	ctx handlers.Context,
	e events.PostEvent,
	next handlers.NextFn[events.PostEvent],
) error {
	logger := ctx.Logger().With(zap.String("handler", "butts"))

	logger.Debug("Asked me to show butts")

	factories.SendTyping(ctx.WSClient(), e.ChannelID())

	butt, err := bh.Butts.Random(ctx.Context())
	if err != nil {
		return fmt.Errorf("failed to get butts: %w", err)
	}

	logger.Debug("Send butts", zap.String("butts", butt.URL))

	post := &model.Post{ChannelId: e.ChannelID(), Message: butt.URL}

	if _, err = ctx.Client().SendPost(ctx.Context(), post); err != nil {
		return fmt.Errorf("failed to answer: %w", err)
	}

	e.Ack()
	return nil
}
