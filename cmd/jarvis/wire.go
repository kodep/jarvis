//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kodep/jarvis/internal/oboobs"
)

func InitializeApp() (App, func(), error) {
	panic(
		wire.Build(
			EventsHandlersSet,
			ProvideApp,
			ProvideConfig,
			ProvideListener,
			ProvideLogger,
			ProvideMattermostClient,
			oboobs.NewBoobsClient,
			oboobs.NewButtsClient,
		),
	)
}
