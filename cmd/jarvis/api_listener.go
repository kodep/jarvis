package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kodep/jarvis/internal/api"
	"github.com/kodep/jarvis/internal/birthday"
	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"github.com/mattermost/mattermost/server/public/model"
	"go.uber.org/zap"
)

type APIListener struct {
	birthday *birthday.Generator
	client   *mattermost.Client
	logger   *zap.Logger
	server   *api.Server
	conf     Config
}

func ProvideAPIListener(
	birthday *birthday.Generator,
	client *mattermost.Client,
	logger *zap.Logger,
	server *api.Server,
	conf Config,
) APIListener {
	return APIListener{
		birthday: birthday,
		client:   client,
		logger:   logger,
		server:   server,
		conf:     conf,
	}
}

func (l *APIListener) Listen(ctx context.Context) error {
	l.InitRoutes()

	if err := l.server.Listen(ctx); err != nil {
		return fmt.Errorf("failed to start HTTP listener: %w", err)
	}

	return nil
}

func (l *APIListener) InitRoutes() {
	l.server.Router().Use(l.AuthGuard)
	l.server.Router().
		HandleFunc("/api/v1/congratulate", l.Congratulate).Methods("POST")
}

func (l *APIListener) Congratulate(w http.ResponseWriter, r *http.Request) {
	var err error

	p := congratulatePayload{}
	if err = p.Read(r); err != nil {
		l.RespondError(w, err, http.StatusBadRequest)
		return
	}

	msg, err := l.birthday.Generate(r.Context(), l.birthday.Prompt(p.Name, p.Description))
	if err != nil {
		l.RespondError(w, fmt.Errorf("failed to generate birthday message: %w", err), http.StatusServiceUnavailable)
		return
	}

	congratulation := "@channel" + " Поздравляем с Днем Рождения" + " @" + p.NickName + " :tada:!"
	if msg != "" {
		congratulation += "\n" + msg
	}

	post := &model.Post{ChannelId: l.conf.BirthdayChannelID, Message: congratulation}
	if _, err = l.client.SendPost(r.Context(), post); err != nil {
		l.RespondError(w, fmt.Errorf("failed to send birthday message: %w", err), http.StatusServiceUnavailable)
		return
	}
}

func (l *APIListener) RespondError(w http.ResponseWriter, err error, status int) {
	l.logger.Error("failed process a request", zap.Error(err))

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	r, err := json.Marshal(errorResponse{Error: err.Error()})
	if err != nil {
		l.logger.Error("failed to marshal error response", zap.Error(err))
		return
	}

	_, _ = w.Write(r)
}

type congratulatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	NickName    string `json:"nick_name"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (p *congratulatePayload) Read(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		return fmt.Errorf("failed to decode request: %w", err)
	}

	if p.Name == "" || p.Description == "" || p.NickName == "" {
		return errors.New("name, description and nick_name are required")
	}

	return nil
}

func (l *APIListener) AuthGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			l.RespondError(w, err, http.StatusInternalServerError)
			return
		}

		h := hmac.New(sha256.New, []byte(l.conf.APISecret))
		h.Write(body)
		hash := base64.StdEncoding.EncodeToString(h.Sum(nil))
		authHash := r.Header.Get("Authorization")

		if hash != authHash {
			l.RespondError(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}
