package handlers

import (
	"context"

	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"go.uber.org/zap"
)

type (
	WSClient interface {
		SendMessage(action string, data map[string]interface{})
	}

	Context interface {
		Context() context.Context
		Logger() *zap.Logger
		Client() *mattermost.Client
		WSClient() WSClient
	}
)

type hContext struct {
	ctx      context.Context
	logger   *zap.Logger
	client   *mattermost.Client
	wsClient *mattermost.WSClient
}

var _ Context = &hContext{}

func NewContext(
	ctx context.Context,
	logger *zap.Logger,
	client *mattermost.Client,
	wsClient *mattermost.WSClient,
) Context {
	return &hContext{
		ctx:      ctx,
		logger:   logger,
		client:   client,
		wsClient: wsClient,
	}
}

func (c *hContext) Context() context.Context {
	return c.ctx
}

func (c *hContext) Logger() *zap.Logger {
	return c.logger
}

func (c *hContext) Client() *mattermost.Client {
	return c.client
}

func (c *hContext) WSClient() WSClient {
	return c.wsClient
}
