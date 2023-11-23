package oboobs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type baseClient struct {
	apiURL string
	cdnURL string
	prefix string
}

type baseEntity struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

func (c baseClient) getRandom(ctx context.Context) (baseEntity, error) {
	url, _ := url.JoinPath(c.apiURL, c.prefix, "/0/1/random") // start=0, limit=1

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return baseEntity{}, fmt.Errorf("failed to create OpenBoobs request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return baseEntity{}, fmt.Errorf("failed to perform OpenBoobs request: %w", err)
	}

	defer res.Body.Close()

	r := []struct {
		ID      int64  `json:"id"`
		Preview string `json:"preview"`
	}{}

	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(&r); err != nil {
		return baseEntity{}, fmt.Errorf("failed to decode OpenBoobs response: %w", err)
	}

	if len(r) == 0 {
		return baseEntity{}, fmt.Errorf("OpenBoobs returned 0 entities")
	}

	return baseEntity{
		ID:  r[0].ID,
		URL: c.getPreviewURL(r[0].Preview),
	}, nil
}

func (c baseClient) getPreviewURL(preview string) string {
	s, _ := url.JoinPath(c.cdnURL, preview)
	return s
}
