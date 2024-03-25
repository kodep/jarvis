package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Options struct {
	Host   string
	Logger *zap.Logger
}

type Client struct {
	client *http.Server
	router *mux.Router
	logger *zap.Logger
}

const WriteTimeout = 15 * time.Second
const ReadTimeout = 15 * time.Second

func NewClient(options Options) *Client {
	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         options.Host,
		WriteTimeout: WriteTimeout,
		ReadTimeout:  ReadTimeout,
	}

	return &Client{
		client: srv,
		logger: options.Logger,
		router: r,
	}
}

func (c *Client) Router() *mux.Router {
	return c.router
}

func (c *Client) ListenAndServe() {
	err := c.client.ListenAndServe()

	if err != nil {
		c.logger.Error("HTTP server start failed: ", zap.Error(err))
	}
}
