package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/enrico5b1b4/tbwrap"
)

func main() {
	telegramBotToken := mustGetEnv("TELEGRAM_BOT_TOKEN")
	allowedChats := parseAllowedChats(os.Getenv("ALLOWED_CHATS"))

	botConfig := tbwrap.Config{
		Token:        telegramBotToken,
		AllowedChats: allowedChats,
	}
	telegramBot, err := tbwrap.NewBot(botConfig)
	if err != nil {
		log.Fatal(err)
	}

	todos := make(map[int64][]string)

	telegramBot.Handle("/todo list", HandleList(todos))
	telegramBot.HandleRegExp(`\/todo add (?P<value>.*)`, HandleAdd(todos))
	telegramBot.HandleRegExp(`\/todo remove (?P<index>.*)`, HandleRemove(todos))
	telegramBot.Start()
}

func mustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln(fmt.Sprintf("%s must be set", name))
	}

	return value
}

func parseAllowedChats(list string) []int {
	if len(list) == 0 {
		return nil
	}

	sepList := strings.Split(list, ",")
	intList := make([]int, len(sepList))
	var err error

	for i := range sepList {
		intList[i], err = strconv.Atoi(strings.TrimSpace(sepList[i]))
		if err != nil {
			log.Fatalln(err)
		}
	}

	return intList
}
