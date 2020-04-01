package tbwrap

//go:generate mockgen -destination=./mocks/mock_Bot.go -package=mocks github.com/enrico5b1b4/tbwrap Bot
//go:generate mockgen -destination=./mocks/mock_TeleBot.go -package=mocks github.com/enrico5b1b4/tbwrap TeleBot

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type HandlerFunc func(c Context) error

type Route struct {
	Path    *regexp.Regexp
	Handler HandlerFunc
}

type TeleBot interface {
	Handle(endpoint interface{}, handler interface{})
	Delete(message tb.Editable) error
	Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error
	Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
	Start()
}

type WrapBot struct {
	tBot   TeleBot
	routes []*Route
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
			routes: []*Route{},
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
		routes: []*Route{},
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

type messageToDelete struct {
	MessageID string
	ChatID    int64
}

// nolint:gocritic
func (m messageToDelete) MessageSig() (string, int64) {
	return m.MessageID, m.ChatID
}

func (b *WrapBot) Delete(chatID int64, messageID int) error {
	msgToDelete := messageToDelete{ChatID: chatID, MessageID: strconv.Itoa(messageID)}

	return b.tBot.Delete(msgToDelete)
}

func (b *WrapBot) HandleRegExp(path string, handler HandlerFunc) {
	compiledRegExp := regexp.MustCompile(path)

	b.routes = append(b.routes, &Route{Path: compiledRegExp, Handler: handler})
}

func (b *WrapBot) HandleMultiRegExp(paths []string, handler HandlerFunc) {
	for i := range paths {
		compiledRegExp := regexp.MustCompile(paths[i])

		b.routes = append(b.routes, &Route{Path: compiledRegExp, Handler: handler})
	}
}

func (b *WrapBot) Handle(path string, handler HandlerFunc) {
	b.handle(path, func(m *tb.Message) {
		c := NewContext(b, m, nil, nil)
		err := handler(c)
		if err != nil {
			_, _ = c.Send(fmt.Sprintf("%s", err))
			log.Println(err)
		}
	})
}

func (b *WrapBot) HandleButton(path *tb.InlineButton, handler HandlerFunc) {
	b.handle(path, func(callback *tb.Callback) {
		c := NewContext(
			b,
			callback.Message,
			callback,
			nil,
		)
		err := handler(c)
		if err != nil {
			_, _ = c.Send(fmt.Sprintf("%s", err))
			log.Println(err)
		}
	})
}

func (b *WrapBot) handleOnText(m *tb.Message) {
	for i := range b.routes {
		matches := b.routes[i].Path.FindStringSubmatch(m.Text)

		if len(matches) > 0 {
			c := NewContext(b, m, nil, b.routes[i].Path)
			err := b.routes[i].Handler(c)
			if err != nil {
				_, _ = c.Send(fmt.Sprintf("%s", err))
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
		b.handleOnText(m)
	})

	b.tBot.Start()
}
