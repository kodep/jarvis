package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Options struct {
	Host   string
	Logger *zap.Logger
}

type ApiClient struct {
	client *http.Server
	router *mux.Router
	logger *zap.Logger
}

func NewClient(options Options) *ApiClient {
	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         options.Host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &ApiClient{
		client: srv,
		logger: options.Logger,
		router: r,
	}

}

func (c *ApiClient) Router() *mux.Router {
	return c.router
}

func (c *ApiClient) ListenAndServe(ctx context.Context) {
	err := c.client.ListenAndServe()

	if err != nil {
		c.logger.Error("HTTP server start failed: ", zap.Error(err))
	}
}
