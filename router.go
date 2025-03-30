package tgool

import "log"

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
		}

		break
	}
}
