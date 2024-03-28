package main

import (
	"context"

	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"github.com/kodep/jarvis/internal/mattermost/events"
	"github.com/kodep/jarvis/internal/mattermost/handlers"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type ChatListener struct {
	client   *mattermost.Client
	wsClient *mattermost.WSClient
	logger   *zap.Logger
	handler  EventsHandler
}

func ProvideChatListener(
	logger *zap.Logger,
	client *mattermost.Client,
	wsClient *mattermost.WSClient,
	handler EventsHandler,
) ChatListener {
	l := ChatListener{
		client:   client,
		wsClient: wsClient,
		logger:   logger,
		handler:  handler,
	}
	return l
}

func (l *ChatListener) Listen(ctx context.Context) {
	const (
		channelSize     = 10
		listenersNumber = 10
	)

	chs := mattermost.WSListenChannels{
		Events: make(chan events.Event, channelSize),
		Errors: make(chan error, channelSize),
	}

	defer close(chs.Events)
	defer close(chs.Errors)

	go l.wsClient.Start(ctx, listenersNumber, chs)

	handleEvent := func(e events.Event) {
		lo.TryCatchWithErrorValue(func() error {
			l.handleEvent(ctx, e)
			return nil
		}, func(val any) {
			l.logger.Error("Event handler panicked", zap.Any("Error", val))
		})
	}

	handleError := func(err error) {
		if eventErr, ok := lo.ErrorsAs[*events.UnknownEventError](err); ok {
			l.logger.Debug("Unknown event", zap.String("EventType", string(eventErr.WsEvent.EventType())))
		} else {
			l.logger.Error("Mattermost WebSocket failed:", zap.Error(err))
		}
	}

	for {
		select {
		case e := <-chs.Events:
			go handleEvent(e)
		case err := <-chs.Errors:
			handleError(err)
		case <-ctx.Done():
			return
		}
	}
}

func (l *ChatListener) handleEvent(ctx context.Context, e events.Event) {
	hCtx := handlers.NewContext(ctx, l.logger, l.client, l.wsClient)

	if err := l.handler(hCtx, e); err != nil {
		l.logger.Warn("Failed to process event", zap.String("EventType", string(e.RawEventType())), zap.Error(err))
	}
}
