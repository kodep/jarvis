package main

import (
	"context"

	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"go.uber.org/zap"
)

type App struct {
	logger   *zap.Logger
	client   *mattermost.Client
	listener Listener
}

func ProvideApp(logger *zap.Logger, client *mattermost.Client, listener Listener) App {
	return App{logger, client, listener}
}

func (a App) Run(ctx context.Context) {
	if err := a.client.Connect(); err != nil {
		a.logger.Fatal("Failed to connect to mattermost", zap.Error(err))
	}

	a.logger.Info("Connected to Mattermost",
		zap.String("ID", a.client.User().Id),
		zap.String("User", a.client.User().Username),
		zap.String("Team", a.client.Team().Name),
	)

	a.logger.Info("Listen for events")
	a.listener.Listen(ctx)

	a.logger.Info("Shutting down")
}
