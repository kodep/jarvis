package handlers

import (
	"github.com/kodep/jarvis/internal/mattermost/events"
)

type NextFn[T events.Event] func(ctx Context, e T) error
type MiddlwareFn[T events.Event] func(ctx Context, e T, next NextFn[T]) error
type HandlerFn[T events.Event] func(ctx Context, e T) error
type FilterFn[T events.Event] func(ctx Context, e T) (bool, error)

func Filter[T events.Event](fn FilterFn[T], h MiddlwareFn[T]) MiddlwareFn[T] {
	return func(ctx Context, e T, next NextFn[T]) error {
		ok, err := fn(ctx, e)
		if err != nil {
			return err
		}

		if ok {
			return h(ctx, e, next)
		}

		return next(ctx, e)
	}
}

func Handle[T events.Event](h HandlerFn[T]) MiddlwareFn[T] {
	return func(ctx Context, e T, next NextFn[T]) error {
		e.Ack()
		return h(ctx, e)
	}
}

func Pipe[T events.Event](handlers ...MiddlwareFn[T]) MiddlwareFn[T] {
	return func(ctx Context, e T, next NextFn[T]) error {
		return runMiddlewares(ctx, e, next, handlers)
	}
}

func runMiddlewares[T events.Event](ctx Context, e T, next NextFn[T], handlers []MiddlwareFn[T]) error {
	select {
	case <-ctx.Context().Done():
		return ctx.Context().Err() //nolint:wrapcheck // don't wrap context.Canceled errors
	default:
	}

	if len(handlers) == 0 {
		return next(ctx, e)
	}

	handler, rest := handlers[0], handlers[1:]

	return handler(ctx, e, func(ctx Context, e T) error {
		return runMiddlewares(ctx, e, next, rest)
	})
}

func TerminatePipe(h MiddlwareFn[events.Event]) func(ctx Context, e events.Event) error {
	return func(ctx Context, e events.Event) error {
		return h(ctx, e, func(Context, events.Event) error {
			return nil
		})
	}
}
