package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BirthdayChannelID string
	BoobsChannelID    string
	BotTeamName       string
	BotToken          string
	HTTPPort          int
	MattermostURL     string
	Mode              ModeType
	OpenAIKey         string
	APISecret         string
}

type Env = string
type ModeType string

const (
	BirthdayChannelID Env = "BIRTHDAY_CHANNEL_ID"
	BoobsChannelID    Env = "BOOBS_CHANNEL_ID"
	BotTeamName       Env = "BOT_TEAM_NAME"
	BotToken          Env = "BOT_TOKEN"
	HTTPPort          Env = "HTTP_PORT"
	MattermostURL     Env = "MATTERMOST_URL"
	OpenAIKey         Env = "OPENAI_KEY"
	APISecret         Env = "API_SECRET"
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

	rHTTPPort := os.Getenv(HTTPPort)
	if rHTTPPort == "" {
		rHTTPPort = "8080"
	}

	httpPort, err := strconv.Atoi(rHTTPPort)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse HTTP port: %w", err)
	}

	c := Config{
		Mode:              mode,
		MattermostURL:     os.Getenv(MattermostURL),
		BotToken:          os.Getenv(BotToken),
		BotTeamName:       os.Getenv(BotTeamName),
		BoobsChannelID:    os.Getenv(BoobsChannelID),
		HTTPPort:          httpPort,
		BirthdayChannelID: os.Getenv(BirthdayChannelID),
		OpenAIKey:         os.Getenv(OpenAIKey),
		APISecret:         os.Getenv(APISecret),
	}

	if err = c.validate(); err != nil {
		return c, err
	}

	return c, nil
}

func (c Config) validate() error {
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
