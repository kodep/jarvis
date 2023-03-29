package oboobs

import (
	"context"
)

type BoobsClient struct {
	baseClient
}

type BoobsEntity struct {
	ID  int64
	URL string
}

func NewBoobsClient() BoobsClient {
	const (
		apiURL = "http://api.oboobs.ru"
		cdnURL = "https://media.oboobs.ru"
		prefix = "boobs"
	)

	return BoobsClient{
		baseClient: baseClient{
			apiURL: apiURL,
			cdnURL: cdnURL,
			prefix: prefix,
		},
	}
}

func (c BoobsClient) Random(ctx context.Context) (BoobsEntity, error) {
	e, err := c.baseClient.getRandom(ctx)
	if err != nil {
		return BoobsEntity{}, err
	}

	return BoobsEntity(e), nil
}
