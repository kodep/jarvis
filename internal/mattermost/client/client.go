package client

import (
	"context"
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
)

type Client struct {
	url      string
	token    string
	teamName string

	client *model.Client4
	team   *model.Team
	user   *model.User
}

type Options struct {
	APIURL   string
	Token    string
	TeamName string
}

func New(options Options) *Client {
	client := model.NewAPIv4Client(options.APIURL)
	client.SetToken(options.Token)

	return &Client{
		url:      options.APIURL,
		token:    options.Token,
		teamName: options.TeamName,
		client:   client,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	if _, _, err := c.client.GetPing(ctx); err != nil {
		return fmt.Errorf("faield to ping server %w", err)
	}

	user, _, err := c.client.GetMe(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get bot account: %w", err)
	}

	team, _, err := c.client.GetTeamByName(ctx, c.teamName, "")
	if err != nil {
		return fmt.Errorf("failed to get bot team: %w", err)
	}

	c.user = user
	c.team = team

	return nil
}

func (c *Client) SendPost(ctx context.Context, post *model.Post) (*model.Post, error) {
	p, _, err := c.client.CreatePost(ctx, post)

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return p, nil
}

func (c *Client) Team() *model.Team {
	return c.team
}

func (c *Client) User() *model.User {
	return c.user
}
