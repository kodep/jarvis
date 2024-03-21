package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MattermostURL     string
	BotToken          string
	BotTeamName       string
	BoobsChannelID    string
	Mode              ModeType
	BirthdayChannelID string
	ApiURL            string
}

type Env = string
type ModeType string

const (
	MattermostURL     Env = "MATTERMOST_URL"
	BotToken          Env = "BOT_TOKEN"
	BotTeamName       Env = "BOT_TEAM_NAME"
	BoobsChannelID    Env = "BOOBS_CHANNEL_ID"
	BirthdayChannelID Env = "BIRTHDAY_CHANNEL_ID"
	ApiURL            Env = "API_URL"
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

	if v := os.Getenv(ApiURL); v != "" {
		c.ApiURL = v
	} else {
		return c, fmt.Errorf("%s environment variable is missing", ApiURL)
	}

	if v := os.Getenv(BirthdayChannelID); v != "" {
		c.BirthdayChannelID = v
	} else {
		return c, fmt.Errorf("%s environment variable is missing", BirthdayChannelID)
	}

	return c, nil
}
