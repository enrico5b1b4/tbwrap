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
	Message() *tb.Message
	ChatID() int64
	Text() string
	Callback() *tb.Callback
	Respond(callback *tb.Callback, response ...*tb.CallbackResponse) error
	Send(msg string, options ...interface{}) (*tb.Message, error)
	Delete(chatID int64, messageID int) error
}

type Bot interface {
	Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error
	Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
	Delete(chatID int64, messageID int) error
}

type context struct {
	bot      Bot
	message  *tb.Message
	callback *tb.Callback
	params   map[string]string
	route    *regexp.Regexp
}

func NewContext(
	bot Bot,
	message *tb.Message,
	callback *tb.Callback,
	route *regexp.Regexp,
) Context {
	var params map[string]string
	if route != nil {
		matches := route.FindStringSubmatch(message.Text)
		names := route.SubexpNames()
		params = mapSubexpNames(matches, names)
	}

	return &context{
		bot:      bot,
		message:  message,
		callback: callback,
		route:    route,
		params:   params,
	}
}

func (c *context) Param(key string) string {
	param := c.params[key]

	return param
}

func (c *context) ChatID() int64 {
	return c.message.Chat.ID
}

func (c *context) Message() *tb.Message {
	return c.message
}

func (c *context) Text() string {
	return c.message.Text
}

func (c *context) Callback() *tb.Callback {
	return c.callback
}

func (c *context) Respond(callback *tb.Callback, response ...*tb.CallbackResponse) error {
	return c.bot.Respond(callback, response...)
}

func (c *context) Send(msg string, options ...interface{}) (*tb.Message, error) {
	sentMsg, err := c.bot.Send(c.message.Chat, msg, options...)

	return sentMsg, err
}

func (c *context) Bind(i interface{}) error {
	return capture.Parse(c.route.String(), c.Text(), i)
}

func (c *context) Delete(chatID int64, messageID int) error {
	return c.bot.Delete(chatID, messageID)
}

func mapSubexpNames(m, n []string) map[string]string {
	if len(m) == 0 || len(n) == 0 {
		return nil
	}

	m, n = m[1:], n[1:]
	r := make(map[string]string, len(m))
	for i := range n {
		r[n[i]] = m[i]
	}
	return r
}
