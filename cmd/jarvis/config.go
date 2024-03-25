package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	MattermostURL     string
	BotToken          string
	BotTeamName       string
	BoobsChannelID    string
	Mode              ModeType
	BirthdayChannelID string
	APIURL            string
}

type Env = string
type ModeType string

const (
	MattermostURL     Env = "MATTERMOST_URL"
	BotToken          Env = "BOT_TOKEN"
	BotTeamName       Env = "BOT_TEAM_NAME"
	BoobsChannelID    Env = "BOOBS_CHANNEL_ID"
	BirthdayChannelID Env = "BIRTHDAY_CHANNEL_ID"
	APIURL            Env = "API_URL"
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

	c := Config{
		Mode:              mode,
		MattermostURL:     os.Getenv(MattermostURL),
		BotToken:          os.Getenv(BotToken),
		BotTeamName:       os.Getenv(BotTeamName),
		BoobsChannelID:    os.Getenv(BoobsChannelID),
		APIURL:            os.Getenv(APIURL),
		BirthdayChannelID: os.Getenv(BirthdayChannelID),
	}

	if err := validate(c); err != nil {
		return c, err
	}

	return c, nil
}

func validate(c Config) error {
	v := reflect.ValueOf(c)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.IsZero() {
			fieldName := reflect.TypeOf(c).Field(i).Name
			return fmt.Errorf("%s environment variable is missing", fieldName)
		}
	}

	return nil
}
