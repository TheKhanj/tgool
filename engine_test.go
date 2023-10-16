package tgool

import (
	"fmt"
	"net/http"
	"testing"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
)

const token = ""

func createHttpClient() (*http.Client, error) {
	full_address := fmt.Sprintf("%s:%d", "127.0.0.1", 10808)
	dialer, err := proxy.SOCKS5("tcp", full_address, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{Dial: dialer.Dial}
	client := &http.Client{Transport: transport}
	return client, nil
}

func getEngine() (*Engine, error) {
	router := NewRouter()
	httpClient, err := createHttpClient()
	if err != nil {
		return nil, err
	}

	tgbot, err := tg.NewBotAPIWithClient(token, tg.APIEndpoint, httpClient)
	if err != nil {
		return nil, err
	}

	return New(3000, router, tgbot), nil
}

func TestEngine(t *testing.T) {
	e, err := getEngine()
	if err != nil {
		panic(err)
	}

	u:=tg.NewUpdate(0)
	u.Timeout = 60

	updates:=e.bot.GetUpdatesChan()
}
