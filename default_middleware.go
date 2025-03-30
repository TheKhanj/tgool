package tgool

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DefaultMiddleWare struct{}

var _ Middleware = &DefaultMiddleWare{}

func (m *DefaultMiddleWare) Handle(
	ctx Context,
	next func(),
) tg.Chattable {
	log.Printf("tgool: default-middleware: no route matched (%s)", ctx.GetRoute())

	return nil
}
