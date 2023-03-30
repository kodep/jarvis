package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MattermostURL  string
	BotToken       string
	BotTeamName    string
	BoobsChannelID string
	Mode           ModeType
}

type Env = string
type ModeType string

const (
	MattermostURL  Env = "MATTERMOST_URL"
	BotToken       Env = "BOT_TOKEN"
	BotTeamName    Env = "BOT_TEAM_NAME"
	BoobsChannelID Env = "BOOBS_CHANNEL_ID"
)

const (
	ModeDevelopment ModeType = "development"
	ModeProduction  ModeType = "production"
)

var DefaultMode = "development" //nolint:gochecknoglobals // It's injected variable

func ProvideConfig() (Config, error) {
	_ = godotenv.Load()

	mode := ModeDevelopment
	if DefaultMode == "production" {
		mode = ModeProduction
	}

	c := Config{Mode: mode}

	if v := os.Getenv(MattermostURL); v != "" {
		c.MattermostURL = v
	} else {
		return c, fmt.Errorf("%s environment variable is missing", MattermostURL)
	}

	if v := os.Getenv(BotToken); v != "" {
		c.BotToken = v
	} else {
		return c, fmt.Errorf("%s environment variable is missing", BotToken)
	}

	if v := os.Getenv(BotTeamName); v != "" {
		c.BotTeamName = v
	} else {
		return c, fmt.Errorf("%s environment variable is missing", BotTeamName)
	}

	if v := os.Getenv(BoobsChannelID); v != "" {
		c.BoobsChannelID = v
	} else {
		return c, fmt.Errorf("%s environment variable is missing", BoobsChannelID)
	}

	return c, nil
}