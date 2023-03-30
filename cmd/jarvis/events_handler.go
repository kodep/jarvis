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
	logsHandler handlers.MiddlwareFn[events.Event]
)

//nolint:gochecknoglobals // it's needed for wire
var EventsHandlersSet = wire.NewSet(
	provideEventsHandler,
	provideLogsHandler,
	wire.Struct(new(boobsAndButtsHandler), "*"),
)

func provideEventsHandler(logsHandler logsHandler, conf Config, bb boobsAndButtsHandler) EventsHandler {
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

func provideLogsHandler(logger *zap.Logger) logsHandler {
	return func(ctx context.Context, e events.Event, next handlers.NextFn[events.Event]) error {
		logger.Debug("Received event", zap.String("EventType", e.RawEventType()))

		err := next(handlers.WithLogger(ctx, logger), e)

		if !e.Acknowledged() && err == nil {
			logger.Debug("Skipped event", zap.String("EventType", e.RawEventType()))
		}

		return err
	}
}

type boobsAndButtsHandler struct {
	Client *mattermost.Client
	Boobs  oboobs.BoobsClient
	Butts  oboobs.ButtsClient
}

func (bh boobsAndButtsHandler) handleBoobs(
	ctx context.Context,
	e events.PostEvent,
	next handlers.NextFn[events.PostEvent],
) error {
	logger := handlers.GetLogger(ctx).With(zap.String("handler", "boobs"))

	logger.Debug("Asked me to show boobs")

	boob, err := bh.Boobs.Random(ctx)
	if err != nil {
		return fmt.Errorf("failed to get boobs: %w", err)
	}

	logger.Debug("Send boobs", zap.String("boobs", boob.URL))

	post := &model.Post{ChannelId: e.ChannelID(), Message: boob.URL}

	if _, err = bh.Client.SendPost(post); err != nil {
		return fmt.Errorf("failed to answer: %w", err)
	}

	e.Ack()
	return nil
}

func (bh boobsAndButtsHandler) handleButts(
	ctx context.Context,
	e events.PostEvent,
	next handlers.NextFn[events.PostEvent],
) error {
	logger := handlers.GetLogger(ctx).With(zap.String("handler", "butts"))

	logger.Debug("Asked me to show butts")

	butt, err := bh.Butts.Random(ctx)
	if err != nil {
		return fmt.Errorf("failed to get butts: %w", err)
	}

	logger.Debug("Send butts", zap.String("butts", butt.URL))

	post := &model.Post{ChannelId: e.ChannelID(), Message: butt.URL}

	if _, err = bh.Client.SendPost(post); err != nil {
		return fmt.Errorf("failed to answer: %w", err)
	}

	e.Ack()
	return nil
}
