package main

import (
	"context"

	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"github.com/kodep/jarvis/internal/mattermost/events"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Listener struct {
	client   *mattermost.Client
	wsClient *mattermost.WSClient
	logger   *zap.Logger
	handler  EventsHandler
}

func ProvideListener(
	logger *zap.Logger,
	client *mattermost.Client,
	wsClient *mattermost.WSClient,
	handler EventsHandler,
) Listener {
	l := Listener{
		client:   client,
		wsClient: wsClient,
		logger:   logger,
		handler:  handler,
	}
	return l
}

func (l *Listener) Listen(ctx context.Context) {
	chs := mattermost.WSListenChannels{
		Events: make(chan *model.WebSocketEvent, 1),
		Errors: make(chan error, 1),
	}

	defer close(chs.Events)
	defer close(chs.Errors)

	go l.wsClient.Start(ctx, chs)

	handleEvent := func(e *model.WebSocketEvent) {
		lo.TryCatchWithErrorValue(func() error {
			l.handleEvent(ctx, e)
			return nil
		}, func(val any) {
			l.logger.Error("Event handler panicked", zap.Any("Error", val))
		})
	}

	for {
		select {
		case e := <-chs.Events:
			go handleEvent(e)
		case err := <-chs.Errors:
			l.logger.Error("Mattermost WebSocket failed:", zap.Error(err))
		case <-ctx.Done():
			return
		}
	}
}

func (l *Listener) handleEvent(ctx context.Context, r *model.WebSocketEvent) {
	e, err := events.NewEvent(r)
	if err != nil {
		if _, ok := lo.ErrorsAs[*events.UnknownEventError](err); ok {
			l.logger.Debug("Unknown event", zap.String("EventType", r.EventType()))
		} else {
			l.logger.Warn("Failed to parse event", zap.String("EventType", r.EventType()), zap.Error(err))
		}

		return
	}

	if err = l.handler(ctx, e); err != nil {
		l.logger.Warn("Failed to process event", zap.String("EventType", r.EventType()), zap.Error(err))
	}
}
