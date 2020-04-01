package main

import (
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"gopkg.in/tucnak/telebot.v2"
)

func NewButtons() map[string]*telebot.InlineButton {
	closeCommandBtn := telebot.InlineButton{
		Unique: "closeCommandBtn",
		Text:   "❌ Close",
	}

	return map[string]*telebot.InlineButton{
		"CloseCommandBtn": &closeCommandBtn,
	}
}

func HandleCloseBtn(
	buttons map[string]*telebot.InlineButton,
) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		messageID, err := strconv.Atoi(c.Callback().Data)
		if err != nil {
			return err
		}

		err = c.Respond(c.Callback())
		if err != nil {
			return err
		}

		err = c.Delete(c.ChatID(), messageID)
		if err != nil {
			return err
		}

		err = c.Delete(c.ChatID(), c.Message().ID)
		if err != nil {
			return err
		}

		return nil
	}
}
