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
	buttons := NewButtons()

	telegramBot.Handle(`/show`, HandleShow(buttons))
	telegramBot.HandleButton(buttons["CloseCommandBtn"], HandleCloseBtn(buttons))
	telegramBot.Start()
}

func mustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln(fmt.Sprintf("%s must be set", name))
	}

	return value
}
