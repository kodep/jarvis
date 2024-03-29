package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Options struct {
	Port   int
	Logger *zap.Logger
}

type Server struct {
	server *http.Server
	router *mux.Router
	logger *zap.Logger
	port   int
}

const (
	writeTimeout = 15 * time.Second
	readTimeout  = 15 * time.Second
)

func NewServer(options Options) *Server {
	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", options.Port),
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	return &Server{
		server: srv,
		logger: options.Logger,
		router: r,
		port:   options.Port,
	}
}

func (s *Server) Router() *mux.Router {
	return s.router
}

func (s *Server) Listen(ctx context.Context) error {
	ch := make(chan error)
	defer close(ch)

	s.logger.Info("Starting HTTP server", zap.Int("port", s.port))

	go func() {
		err := s.server.ListenAndServe()

		s.logger.Info("HTTP server stopped")

		if err != http.ErrServerClosed {
			ch <- fmt.Errorf("failed to start HTTP server: %w", err)
		} else {
			ch <- nil
		}
	}()

	go func() {
		<-ctx.Done()
		_ = s.server.Shutdown(ctx)
	}()

	return <-ch
}
