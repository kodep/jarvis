package main

import (
	"context"
	"fmt"
	"sync"

	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"github.com/kodep/jarvis/internal/mattermost/events"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Listener struct {
	client   *mattermost.Client
	wsClient *model.WebSocketClient
	logger   *zap.Logger
	handler  EventsHandler
}

func ProvideListener(logger *zap.Logger, client *mattermost.Client, handler EventsHandler) (Listener, func()) {
	l := Listener{
		client:  client,
		logger:  logger,
		handler: handler,
	}
	return l, l.Close
}

func (l *Listener) Connect() error {
	var err error

	if l.wsClient, err = l.client.Websocket(); err != nil {
		return fmt.Errorf("failed to create websocket client: %w", err)
	}

	if connectErr := l.wsClient.Connect(); connectErr != nil {
		return fmt.Errorf("failed to connect to websocket: %w", connectErr)
	}

	return nil
}

func (l *Listener) Listen(ctx context.Context) {
	const (
		routines = 3
	)

	var wg sync.WaitGroup

	l.wsClient.Listen()

	// Listen WebSocketEvent, ResponseChannel, PingTimeoutChannel
	// to prevent deadlocks
	//nolint:lll // link
	// https://github.com/mattermost/mattermost-server/blob/529ab959e26b7e5507bd345f271bea9b86219432/model/websocket_client.go#L42
	wg.Add(routines)

	go func() {
		defer wg.Done()
		l.logger.Debug("Listening for WS events")
		l.listenEvents(ctx, l.wsClient.EventChannel)
	}()

	go func() {
		defer wg.Done()
		l.logger.Debug("Listening for WS responses")
		l.listenResponses(ctx, l.wsClient.ResponseChannel)
	}()

	go func() {
		defer wg.Done()
		l.logger.Debug("Listening for WS ping timeouts")
		l.listenPingTimeout(ctx, l.wsClient.PingTimeoutChannel)
	}()

	wg.Wait()
}

func (l *Listener) listenEvents(ctx context.Context, ch <-chan *model.WebSocketEvent) {
	for {
		select {
		case r := <-ch:
			if r == nil {
				continue
			}

			go func(e *model.WebSocketEvent) {
				lo.TryCatchWithErrorValue(func() error {
					l.handleEvent(ctx, e)
					return nil
				}, func(val any) {
					l.logger.Error("Event handler panicked", zap.Any("Error", val))
				})
			}(r)
		case <-ctx.Done():
			return
		}
	}
}

func (l *Listener) listenResponses(ctx context.Context, ch <-chan *model.WebSocketResponse) {
	for {
		select {
		case <-ch:
		case <-ctx.Done():
			return
		}
	}
}

func (l *Listener) listenPingTimeout(ctx context.Context, ch <-chan bool) {
	for {
		select {
		case <-ch:
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

func (l Listener) Close() {
	if l.wsClient != nil {
		l.wsClient.Close()
	}
}
