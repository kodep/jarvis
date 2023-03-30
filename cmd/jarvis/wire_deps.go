package main

import (
	"fmt"

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
