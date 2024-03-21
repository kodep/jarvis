package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kodep/jarvis/internal/api"
	"github.com/kodep/jarvis/internal/birthday"
	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"github.com/mattermost/mattermost/server/public/model"
	"go.uber.org/zap"
)

type APIListener struct {
	apiClient *api.Client
	client    *mattermost.Client
	logger    *zap.Logger
	conf      Config
}

func ProvideAPIListener(
	apiClient *api.Client,
	client *mattermost.Client,
	logger *zap.Logger,
	conf Config,
) APIListener {
	return APIListener{
		apiClient: apiClient,
		client:    client,
		logger:    logger,
		conf:      conf,
	}
}

func (l *APIListener) ListenAPI(ctx context.Context) {
	l.InitRoutes(ctx)
	go l.apiClient.ListenAndServe()
}

func (l *APIListener) InitRoutes(ctx context.Context) {
	r := l.apiClient.Router()

	r.HandleFunc("/api/v1/congratulate", func(w http.ResponseWriter, r *http.Request) {
		l.Congratulate(ctx, w, r)
	}).Methods("POST")
}

func (l *APIListener) Congratulate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	message, err := birthday.GetMessage(r)
	if err != nil {
		l.logger.Error("getting birthday message error", zap.Error(err))
		l.SendError(w, err, "Getting birthday message error")
		return
	}

	post := &model.Post{ChannelId: l.conf.BirthdayChannelID, Message: message}

	if _, err = l.client.SendPost(ctx, post); err != nil {
		l.logger.Error("sending post to client failed", zap.Error(err))
		l.SendError(w, err, "Sending post to client failed")
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"data": message})
	if err != nil {
		l.logger.Error("encoding JSON error: %w", zap.Error(err))
		l.SendError(w, err, "Encoding JSON error")
	}
}

func (l *APIListener) SendError(w http.ResponseWriter, err error, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	encodingErr := json.NewEncoder(w).Encode(map[string]string{"message": msg, "error": err.Error()})

	if encodingErr != nil {
		l.logger.Error("encoding JSON error: %w", zap.Error(err))
	}
}
