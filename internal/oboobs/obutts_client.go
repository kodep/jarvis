package oboobs

import (
	"context"
)

type ButtsClient struct {
	baseClient
}

type ButtsEntity struct {
	ID  int64
	URL string
}

func NewButtsClient() *ButtsClient {
	const (
		apiURL = "http://api.obutts.ru"
		cdnURL = "https://media.obutts.ru"
		prefix = "butts"
	)

	return &ButtsClient{
		baseClient: baseClient{
			apiURL: apiURL,
			cdnURL: cdnURL,
			prefix: prefix,
		},
	}
}

func (c ButtsClient) Random(ctx context.Context) (ButtsEntity, error) {
	e, err := c.baseClient.getRandom(ctx)
	if err != nil {
		return ButtsEntity{}, err
	}

	return ButtsEntity(e), nil
}
