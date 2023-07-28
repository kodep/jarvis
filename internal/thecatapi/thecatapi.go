package thecatapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct{}

type Entity struct {
	ID  string
	URL string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetCat(ctx context.Context) (Entity, error) {
	const url = "https://api.thecatapi.com/v1/images/search"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to create TheCatAPI request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to perform TheCatAPI request: %w", err)
	}

	defer res.Body.Close()

	r := []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}{}

	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(&r); err != nil {
		return Entity{}, fmt.Errorf("failed to decode TheCatAPI response: %w", err)
	}

	if len(r) == 0 {
		return Entity{}, fmt.Errorf("TheCatAPI returned no results")
	}

	return Entity{ID: r[0].ID, URL: r[0].URL}, nil
}
