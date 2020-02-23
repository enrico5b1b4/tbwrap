package tbwrap

//go:generate mockgen -destination=./mocks/mock_Context.go -package=mocks github.com/enrico5b1b4/tbwrap Context

import (
	"regexp"

	"github.com/enrico5b1b4/capture"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Context interface {
	Param(key string) string
	Bind(i interface{}) error
	ChatID() int
	Text() string
	Callback() *tb.Callback
	Respond(callback *tb.Callback, response ...*tb.CallbackResponse) error
	Send(msg string, options ...interface{}) error
}

type context struct {
	bot      Bot
	chat     *tb.Chat
	text     string
	callback *tb.Callback
	chatID   int
	params   map[string]string
	route    *regexp.Regexp
}

func (c *context) Param(key string) string {
	param := c.params[key]

	return param
}

func (c *context) ChatID() int {
	return c.chatID
}

func (c *context) Text() string {
	return c.text
}

func (c *context) Callback() *tb.Callback {
	return c.callback
}

func (c *context) Respond(callback *tb.Callback, response ...*tb.CallbackResponse) error {
	return c.bot.Respond(callback, response...)
}

func (c *context) Send(msg string, options ...interface{}) error {
	_, err := c.bot.Send(c.chat, msg, options...)

	return err
}

func (c *context) Bind(i interface{}) error {
	return capture.Parse(c.route.String(), c.Text(), i)
}
