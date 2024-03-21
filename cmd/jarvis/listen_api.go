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

type ApiListener struct {
	apiClient *api.ApiClient
	client    *mattermost.Client
	logger    *zap.Logger
	conf      Config
}

func ProvideApiListener(
	apiClient *api.ApiClient,
	client *mattermost.Client,
	logger *zap.Logger,
	conf Config,
) ApiListener {
	return ApiListener{
		apiClient: apiClient,
		client:    client,
		logger:    logger,
		conf:      conf,
	}
}

func (l *ApiListener) ListenApi(ctx context.Context) {
	l.InitRoutes(ctx)
	go l.apiClient.ListenAndServe(ctx)
}

func (l *ApiListener) InitRoutes(ctx context.Context) {
	r := l.apiClient.Router()

	r.HandleFunc("/api/v1/congratulate", func(w http.ResponseWriter, r *http.Request) { l.Congratulate(ctx, w, r) }).Methods("POST")
}

func (l *ApiListener) Congratulate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	message, err := birthday.GetMessage(r)
	if err != nil {
		l.logger.Error("getting birthday message error", zap.Error(err))
		SendError(w, err, "Getting birthday message error")
		return
	}

	post := &model.Post{ChannelId: l.conf.BirthdayChannelID, Message: message}

	if _, err = l.client.SendPost(ctx, post); err != nil {
		l.logger.Error("sending post to client failed", zap.Error(err))
		SendError(w, err, "Sending post to client failed")
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"data": message})
}

func SendError(w http.ResponseWriter, err error, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(map[string]string{"message": msg, "error": err.Error()})
}
