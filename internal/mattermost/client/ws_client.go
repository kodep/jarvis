package client

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

const (
	wsRetryTimeout     = 10 * time.Second
	wsReconnectTimeout = 10 * time.Second
	wsListenersCount   = 4
)

type wsDisonnected bool

type WSClientOptions struct {
	URL    string
	Token  string
	Logger *zap.Logger
}

type WSClient struct {
	url   string
	token string

	ws     *model.WebSocketClient
	logger *zap.Logger
}

type WSListenChannels struct {
	Events chan *model.WebSocketEvent
	Errors chan error
}

func NewWSClient(options WSClientOptions) *WSClient {
	if options.Logger == nil {
		options.Logger = zap.NewNop()
	}

	return &WSClient{
		url:    options.URL,
		token:  options.Token,
		logger: options.Logger,
	}
}

func GetWSURL(apiURL string) (string, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to construct websocket url: %w", err)
	}

	scheme := u.Scheme

	u.Scheme = "ws"
	if scheme == "https" {
		u.Scheme = "wss"
	}

	return u.String(), nil
}

func (c *WSClient) Start(ctx context.Context, channels WSListenChannels) {
	defer func() {
		if c.ws != nil {
			c.ws.Close()
			c.ws = nil
		}
	}()

	for {
		c.logger.Info("Connect to Mattermost via Websocket")
		c.reconnectUntilReady(ctx, channels.Errors)
		c.logger.Info("Successfully connected to Mattermost via Websocket")

		select {
		case <-ctx.Done():
			return
		default:
		}

		c.logger.Debug("Listen Mattermost events")
		disconnected := c.startListeners(ctx, channels.Events)
		if disconnected {
			c.logger.Debug("Disconneted from Mattermost server. Retrying.")
		}

		select {
		case <-time.After(wsReconnectTimeout):
		case <-ctx.Done():
			return
		}
	}
}

func (c *WSClient) startListeners(ctx context.Context, evCh chan<- *model.WebSocketEvent) wsDisonnected {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resCh := make(chan wsDisonnected, wsListenersCount)

	for i := 0; i < wsListenersCount; i++ {
		go func() {
			r := c.startListener(ctx, evCh)
			resCh <- r
		}()
	}

	c.ws.Listen()

	return <-resCh
}

func (c *WSClient) startListener(ctx context.Context, evCh chan<- *model.WebSocketEvent) wsDisonnected {
	for {
		select {
		case e, ok := <-c.ws.EventChannel:
			if !ok {
				return true
			}
			evCh <- e
		case _, ok := <-c.ws.ResponseChannel:
			if !ok {
				return true
			}
		case _, ok := <-c.ws.PingTimeoutChannel:
			if !ok {
				return true
			}
		case <-ctx.Done():
			c.ws.Close()
			c.ws = nil
			return false
		}
	}
}

func (c *WSClient) reconnectUntilReady(ctx context.Context, errCh chan<- error) {
	for {
		err := c.connect()

		if err == nil {
			return
		}

		c.logger.Debug("Unable to connect to Mattermost. Retrying.", zap.Error(err))

		select {
		case errCh <- err:
		default:
		}

		select {
		case <-time.After(wsRetryTimeout):
		case <-ctx.Done():
			return
		}
		continue
	}
}

func (c *WSClient) connect() error {
	var err error

	if c.ws != nil {
		c.ws.Close()
		c.ws = nil
	}

	c.ws, err = model.NewWebSocketClient(c.url, c.token)
	if err != nil {
		return fmt.Errorf("failed to create websocket client: %w", err)
	}

	if aerr := c.ws.Connect(); aerr != nil {
		return fmt.Errorf("failed to connect to Mattermost: %w", aerr)
	}

	return nil
}
