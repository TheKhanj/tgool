package tgool

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Engine struct {
	router *Router
	token  string
	port   int
	bot    *tg.BotAPI
}

func NewEngine(
	router *Router,
	bot *tg.BotAPI,
) *Engine {
	return &Engine{
		router: router,
		bot:    bot,
	}
}

func (e *Engine) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	log.Println("tgool: started handling updates")
	defer log.Println("tgool: stopped handling updates")

	chatsState := &chatsState{}

	for update := range updates {
		ctx := context{
			bot:        e.bot,
			u:          &update,
			chatsState: chatsState,
		}
		e.router.Route(&ctx)
	}
}
