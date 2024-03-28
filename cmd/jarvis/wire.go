//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kodep/jarvis/internal/oboobs"
	"github.com/kodep/jarvis/internal/thecatapi"
)

func InitializeApp() (App, func(), error) {
	panic(
		wire.Build(
			EventsHandlersSet,
			oboobs.NewBoobsClient,
			oboobs.NewButtsClient,
			ProvideAPIListener,
			ProvideAPIServer,
			ProvideApp,
			ProvideBirthdayGenerator,
			ProvideChatListener,
			ProvideConfig,
			ProvideLogger,
			ProvideMattermostClient,
			ProvideMattermostWSClient,
			thecatapi.NewClient,
		),
	)
}
