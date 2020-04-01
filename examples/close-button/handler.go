package main

import (
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"gopkg.in/tucnak/telebot.v2"
)

func HandleShow(buttons map[string]*telebot.InlineButton) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		inlineKeys := [][]telebot.InlineButton{
			{*buttons["CloseCommandBtn"]},
		}
		inlineKeys[0][0].Data = strconv.Itoa(c.Message().ID)

		_, err := c.Send("Click the button to close this message", &telebot.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})

		return err
	}
}
