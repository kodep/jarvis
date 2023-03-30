package client

import (
	"fmt"
	"net/url"

	"github.com/mattermost/mattermost-server/v6/model"
)

type Client struct {
	client  *model.Client4
	options Options
	team    *model.Team
	user    *model.User
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
		client:  client,
		options: options,
	}
}

func (c *Client) Connect() error {
	user, _, err := c.client.GetMe("")
	if err != nil {
		return fmt.Errorf("failed to get bot account: %w", err)
	}

	team, _, err := c.client.GetTeamByName(c.options.TeamName, "")
	if err != nil {
		return fmt.Errorf("failed to get bot team: %w", err)
	}

	if _, _, err = c.client.GetOldClientLicense(""); err != nil {
		return fmt.Errorf("faield to get config %w", err)
	}

	c.user = user
	c.team = team

	return nil
}

func (c *Client) SendPost(post *model.Post) (*model.Post, error) {
	p, _, err := c.client.CreatePost(post)

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

func (c *Client) Websocket() (*model.WebSocketClient, error) {
	u, err := url.Parse(c.options.APIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to construct websocket url: %w", err)
	}

	scheme := u.Scheme

	u.Scheme = "ws"
	if scheme == "https" {
		u.Scheme = "wss"
	}

	ws, err := model.NewWebSocketClient4(u.String(), c.options.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %w", err)
	}

	return ws, nil
}
