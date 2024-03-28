package main

import (
	"fmt"

	"github.com/kodep/jarvis/internal/api"
	"github.com/kodep/jarvis/internal/birthday"
	mattermost "github.com/kodep/jarvis/internal/mattermost/client"
	"go.uber.org/zap"
)

func ProvideLogger(config Config) (*zap.Logger, func(), error) {
	var (
		logger *zap.Logger
		err    error
	)

	if config.Mode == ModeDevelopment {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	return logger, func() {
		_ = logger.Sync()
	}, nil
}

func ProvideMattermostClient(conf Config) *mattermost.Client {
	return mattermost.New(mattermost.Options{
		APIURL:   conf.MattermostURL,
		TeamName: conf.BotTeamName,
		Token:    conf.BotToken,
	})
}

func ProvideMattermostWSClient(logger *zap.Logger, conf Config) (*mattermost.WSClient, error) {
	wsURL, err := mattermost.GetWSURL(conf.MattermostURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get websocket URL: %w", err)
	}

	return mattermost.NewWSClient(mattermost.WSClientOptions{
		URL:    wsURL,
		Token:  conf.BotToken,
		Logger: logger,
	}), nil
}

func ProvideAPIServer(logger *zap.Logger, conf Config) *api.Server {
	return api.NewServer(api.Options{
		Port:   conf.HTTPPort,
		Logger: logger,
	})
}

func ProvideBirthdayGenerator(conf Config) *birthday.Generator {
	return birthday.NewGenerator(birthday.Options{
		OpenAIKey: conf.OpenAIKey,
	})
}
