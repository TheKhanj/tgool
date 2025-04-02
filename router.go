package tgool

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	middlewares []Middleware
}

func NewRouter(
	middlewares ...Middleware,
) *Router {
	return &Router{
		middlewares: append(middlewares, &DefaultMiddleWare{}),
	}
}

func (r *Router) Route(ctx Context) {
	log.Println("tgool: handling new message")

	for _, middleware := range r.middlewares {
		nextCalled := false

		next := func() {
			nextCalled = true
		}

		res := middleware.Handle(ctx, next)

		if nextCalled {
			continue
		}

		if res != nil {
			ctx.Bot().Send(res)

			if ctx.Update().CallbackQuery != nil {
				ctx.Bot().Request(
					tgbotapi.NewCallback(ctx.Update().CallbackQuery.ID, ""),
				)
			}
		}

		break
	}
}
