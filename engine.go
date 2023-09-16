package tgool

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	router *Router
	token  string
	port   int
	bot    *tg.BotAPI
}

func New(
	port int,
	router *Router,
	bot *tg.BotAPI,
) *Engine {
	return &Engine{
		router: router,
		port:   port,
		bot:    bot,
	}
}

func (e *Engine) Listen() {
	var wg sync.WaitGroup
	wg.Add(1)

	addr := fmt.Sprintf("0.0.0.0:%d", e.port)

	updates := make(chan tg.Update, e.bot.Buffer)

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		var update tg.Update
		json.Unmarshal(bytes, &update)

		updates <- update
	})

	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		defer wg.Done()
		msg := fmt.Sprintf("admin telegram server is listening on port %d", e.port)
		logrus.Info(msg)
		server.ListenAndServe()
	}()

	chatsState := new(chatsState)
	for update := range updates {
		ctx := context{
			bot:        e.bot,
			u:          &update,
			chatsState: chatsState,
		}
		e.router.Route(&ctx)
	}

	wg.Wait()
}
