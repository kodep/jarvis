package main

import (
	"context"
	"time"

	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"go.uber.org/zap"
)

const (
	connectAttemptTimeout = 10 * time.Second
)

type App struct {
	logger      *zap.Logger
	client      *mattermost.Client
	listener    Listener
	apiListener ApiListener
}

func ProvideApp(logger *zap.Logger, client *mattermost.Client, listener Listener, apiListener ApiListener) App {
	return App{logger, client, listener, apiListener}
}

func (a App) Run(ctx context.Context) {
	a.connectUntilReady(ctx)

	a.logger.Info("Connected to Mattermost",
		zap.String("ID", a.client.User().Id),
		zap.String("User", a.client.User().Username),
		zap.String("Team", a.client.Team().Name),
	)

	a.logger.Info("Listen for API events")
	a.apiListener.ListenApi(ctx)

	a.logger.Info("Listen for events")
	a.listener.Listen(ctx)

	a.logger.Info("Shutting down")
}

func (a App) connectUntilReady(ctx context.Context) {
	var err error

	for {
		if err = a.client.Connect(ctx); err == nil {
			return
		}

		a.logger.Error("Failed to connect to mattermost. Retry", zap.Error(err))

		select {
		case <-time.After(connectAttemptTimeout):
		case <-ctx.Done():
			return
		}
	}
}
