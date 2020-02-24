package tbwrap

//go:generate mockgen -destination=./mocks/mock_Bot.go -package=mocks github.com/enrico5b1b4/tbwrap Bot
//go:generate mockgen -destination=./mocks/mock_TeleBot.go -package=mocks github.com/enrico5b1b4/tbwrap TeleBot

import (
	"fmt"
	"log"

	"time"

	"regexp"

	tb "gopkg.in/tucnak/telebot.v2"
)

type HandlerFunc func(c Context) error

type Route struct {
	Path    *regexp.Regexp
	Handler HandlerFunc
}

type TeleBot interface {
	Handle(endpoint interface{}, handler interface{})
	Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error
	Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
	Start()
}

type Bot interface {
	Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error
	Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
}

type WrapBot struct {
	tBot   TeleBot
	routes map[*regexp.Regexp]*Route
}

type Config struct {
	Token        string
	AllowedChats []int
	TBot         TeleBot
}

var pollerTimeout = 15 * time.Second

func NewBot(cfg Config) (*WrapBot, error) {
	if cfg.TBot != nil {
		return &WrapBot{
			tBot:   cfg.TBot,
			routes: map[*regexp.Regexp]*Route{},
		}, nil
	}

	poller := NewPollerWithAllowedChats(pollerTimeout, cfg.AllowedChats)
	tBot, err := tb.NewBot(tb.Settings{
		Token:  cfg.Token,
		Poller: poller,
	})
	if err != nil {
		return nil, err
	}

	return &WrapBot{
		tBot:   tBot,
		routes: map[*regexp.Regexp]*Route{},
	}, nil
}

func (b *WrapBot) handle(endpoint, handler interface{}) {
	b.tBot.Handle(endpoint, handler)
}

func (b *WrapBot) Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error {
	return b.tBot.Respond(callback, responseOptional...)
}

func (b *WrapBot) Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error) {
	mergedOptions := append([]interface{}{&tb.SendOptions{ParseMode: tb.ModeMarkdown}}, options...)

	return b.tBot.Send(to, what, mergedOptions...)
}

func (b *WrapBot) HandleRegExp(path string, handler HandlerFunc) {
	compiledRegExp := regexp.MustCompile(path)

	b.routes[compiledRegExp] = &Route{Path: compiledRegExp, Handler: handler}
}

func (b *WrapBot) HandleMultiRegExp(paths []string, handler HandlerFunc) {
	for i := range paths {
		compiledRegExp := regexp.MustCompile(paths[i])

		b.routes[compiledRegExp] = &Route{Path: compiledRegExp, Handler: handler}
	}
}

func (b *WrapBot) Handle(path string, handler HandlerFunc) {
	b.handle(path, func(m *tb.Message) {
		c := &context{chat: m.Chat, text: m.Text, chatID: int(m.Chat.ID), bot: b}
		err := handler(c)
		if err != nil {
			_ = c.Send(fmt.Sprintf("%s", err))
			log.Println(err)
		}
	})
}

func (b *WrapBot) HandleButton(path *tb.InlineButton, handler HandlerFunc) {
	b.handle(path, func(callback *tb.Callback) {
		c := &context{
			chat:     callback.Message.Chat,
			text:     callback.Message.Text,
			callback: callback,
			chatID:   int(callback.Message.Chat.ID),
			bot:      b,
		}
		err := handler(c)
		if err != nil {
			_ = c.Send(fmt.Sprintf("%s", err))
			log.Println(err)
		}
	})
}

func (b *WrapBot) handleOnText(text string, chat *tb.Chat) {
	for regExpKey := range b.routes {
		matches := regExpKey.FindStringSubmatch(text)
		names := regExpKey.SubexpNames()

		if len(matches) > 0 {
			params := mapSubexpNames(matches, names)
			c := &context{chat: chat, text: text, params: params, chatID: int(chat.ID), route: regExpKey, bot: b}
			err := b.routes[regExpKey].Handler(c)
			if err != nil {
				_ = c.Send(fmt.Sprintf("%s", err))
				log.Println(err)
			}

			return
		}
	}
}

func (b *WrapBot) Start() {
	b.handle(tb.OnText, func(m *tb.Message) {
		// all the text messages that weren't
		// captured by existing handlers
		b.handleOnText(m.Text, m.Chat)
	})

	b.tBot.Start()
}

func mapSubexpNames(m, n []string) map[string]string {
	m, n = m[1:], n[1:]
	r := make(map[string]string, len(m))
	for i := range n {
		r[n[i]] = m[i]
	}
	return r
}
