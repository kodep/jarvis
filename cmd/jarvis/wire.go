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
			ProvideApp,
			ProvideLogger,
			ProvideConfig,
			ProvideMattermostClient,
			ProvideListener,
			EventsHandlersSet,
			oboobs.NewBoobsClient,
			oboobs.NewButtsClient,
		),
	)
}
