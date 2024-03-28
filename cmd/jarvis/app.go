package main

import (
	"context"
	"fmt"
	"time"

	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	connectAttemptTimeout = 10 * time.Second
)

type App struct {
	logger       *zap.Logger
	client       *mattermost.Client
	chatListener ChatListener
	apiListener  APIListener
}

func ProvideApp(
	logger *zap.Logger,
	client *mattermost.Client,
	chatListener ChatListener,
	apiListener APIListener,
) App {
	return App{logger, client, chatListener, apiListener}
}

func (a App) Run(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		a.connectUntilReady(ctx)

		a.logger.Info("Connected to Mattermost",
			zap.String("ID", a.client.User().Id),
			zap.String("User", a.client.User().Username),
			zap.String("Team", a.client.Team().Name),
		)

		a.logger.Info("Listen for events")
		a.chatListener.Listen(ctx)

		return nil
	})

	g.Go(func() error {
		a.logger.Info("Listen for API events")
		if err := a.apiListener.Listen(ctx); err != nil {
			return fmt.Errorf("failed to start API listener: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		a.logger.Error("App failed", zap.Error(err))
	}

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
