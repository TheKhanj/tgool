package tgool

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler func(ctx Context) tg.Chattable

type Middleware interface {
	Handle(ctx Context, next func()) tg.Chattable
}
