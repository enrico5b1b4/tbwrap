package main

import (
	"fmt"
	"log"
	"os"

	"github.com/enrico5b1b4/tbwrap"
)

func main() {
	telegramBotToken := mustGetEnv("TELEGRAM_BOT_TOKEN")

	botConfig := tbwrap.Config{
		Token: telegramBotToken,
	}
	telegramBot, err := tbwrap.NewBot(botConfig)
	if err != nil {
		log.Fatal(err)
	}

	telegramBot.Handle(`/ping`, func(c tbwrap.Context) error {
		_, err := c.Send("pong!")

		return err
	})
	telegramBot.Start()
}

func mustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln(fmt.Sprintf("%s must be set", name))
	}

	return value
}
