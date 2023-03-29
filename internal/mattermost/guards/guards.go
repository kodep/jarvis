package guards

import (
	"context"

	"github.com/kodep/jarvis/internal/mattermost/events"
	"github.com/kodep/jarvis/internal/mattermost/handlers"
)

func EventGuard[T events.Event](h handlers.MiddlwareFn[T]) handlers.MiddlwareFn[events.Event] {
	return func(ctx context.Context, e events.Event, next handlers.NextFn[events.Event]) error {
		if evt, ok := e.(T); ok {
			return h(ctx, evt, func(ctx context.Context, e T) error {
				return next(ctx, e)
			})
		}

		return next(ctx, e)
	}
}
