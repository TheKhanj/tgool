package tgool

import (
	"github.com/sirupsen/logrus"
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
	logrus.Info("[telegram] handling message")

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
