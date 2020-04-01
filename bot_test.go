package tbwrap_test

import (
	"errors"
	"log"
	"testing"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/tbwrap/fakes"
	"github.com/stretchr/testify/assert"
)

func Test_TBWrap_HandleSuccess(t *testing.T) {
	fakeTeleBot := fakes.NewTeleBot()
	wrapBot := NewWrapBot(fakeTeleBot)
	wrapBot.Handle("/test", func(c tbwrap.Context) error {
		assert.Equal(t, "/test", c.Text())
		assert.Equal(t, int64(1), c.ChatID())

		return nil
	})
	wrapBot.Start()

	fakeTeleBot.SimulateIncomingMessageToChat(1, "/test")
}

func Test_TBWrap_HandleError(t *testing.T) {
	fakeTeleBot := fakes.NewTeleBot()
	wrapBot := NewWrapBot(fakeTeleBot)
	wrapBot.Handle("/test", func(c tbwrap.Context) error {
		return errors.New("error")
	})
	wrapBot.Start()

	fakeTeleBot.SimulateIncomingMessageToChat(int64(1), "/test")

	assert.Contains(t, fakeTeleBot.OutboundSendMessages, "error")
}

func Test_TBWrap_HandleRegExpSuccess(t *testing.T) {
	fakeTeleBot := fakes.NewTeleBot()
	wrapBot := NewWrapBot(fakeTeleBot)
	wrapBot.HandleRegExp(`\/remind (?P<who>\w+)`, func(c tbwrap.Context) error {
		assert.Equal(t, "/remind Bob", c.Text())
		assert.Equal(t, "Bob", c.Param("who"))
		assert.Equal(t, int64(1), c.ChatID())

		return nil
	})
	wrapBot.Start()

	fakeTeleBot.SimulateIncomingMessageToChat(1, "/remind Bob")
}

func Test_TBWrap_HandleRegExpError(t *testing.T) {
	fakeTeleBot := fakes.NewTeleBot()
	wrapBot := NewWrapBot(fakeTeleBot)
	wrapBot.HandleRegExp(`\/remind (?P<who>\w+)`, func(c tbwrap.Context) error {
		return errors.New("error")
	})
	wrapBot.Start()

	fakeTeleBot.SimulateIncomingMessageToChat(1, "/remind Bob")

	assert.Contains(t, fakeTeleBot.OutboundSendMessages, "error")
}

func Test_TBWrap_HandleMultiRegExpSuccess(t *testing.T) {
	fakeTeleBot := fakes.NewTeleBot()
	wrapBot := NewWrapBot(fakeTeleBot)
	wrapBot.HandleMultiRegExp([]string{`\/remind (?P<who>\w+)`, `\/tell (?P<who>\w+)`}, func(c tbwrap.Context) error {
		assert.Equal(t, "Bob", c.Param("who"))
		assert.Equal(t, int64(1), c.ChatID())

		return nil
	})
	wrapBot.Start()

	fakeTeleBot.SimulateIncomingMessageToChat(1, "/remind Bob")
}

func Test_TBWrap_HandleMultiRegExpError(t *testing.T) {
	fakeTeleBot := fakes.NewTeleBot()
	wrapBot := NewWrapBot(fakeTeleBot)
	wrapBot.HandleMultiRegExp([]string{`\/remind (?P<who>\w+)`, `\/tell (?P<who>\w+)`}, func(c tbwrap.Context) error {
		return errors.New("error")
	})
	wrapBot.Start()

	fakeTeleBot.SimulateIncomingMessageToChat(1, "/remind Bob")

	assert.Contains(t, fakeTeleBot.OutboundSendMessages, "error")
}

func Test_TBWrap_BindMessage(t *testing.T) {
	type Message struct {
		Who string `regexpGroup:"who"`
	}

	fakeTeleBot := fakes.NewTeleBot()
	wrapBot := NewWrapBot(fakeTeleBot)
	wrapBot.HandleMultiRegExp([]string{`\/remind (?P<who>\w+)`, `\/tell (?P<who>\w+)`}, func(c tbwrap.Context) error {
		m := new(Message)
		err := c.Bind(m)
		assert.NoError(t, err)
		assert.Equal(t, "Bob", m.Who)

		return nil
	})
	wrapBot.Start()

	fakeTeleBot.SimulateIncomingMessageToChat(1, "/remind Bob")
}

func Test_TBWrap_SendMessageFromHandler(t *testing.T) {
	fakeTeleBot := fakes.NewTeleBot()
	wrapBot := NewWrapBot(fakeTeleBot)
	wrapBot.Handle("/test", func(c tbwrap.Context) error {
		_, err := c.Send("a message")

		return err
	})
	wrapBot.Start()

	fakeTeleBot.SimulateIncomingMessageToChat(1, "/test")

	assert.Contains(t, fakeTeleBot.OutboundSendMessages, "a message")
}

func Test_TBWrap_OnlyCallsFirstMatchedHandler(t *testing.T) {
	fakeTeleBot := fakes.NewTeleBot()
	wrapBot := NewWrapBot(fakeTeleBot)
	wrapBot.HandleRegExp(`\/test (?P<who>\w+)`, func(c tbwrap.Context) error {
		_, err := c.Send(`I matched with \/test (?P<who>\w+)`)

		return err
	})
	wrapBot.HandleRegExp(`\/test me`, func(c tbwrap.Context) error {
		_, err := c.Send(`I matched with \/test me`)

		return err
	})
	wrapBot.Start()

	fakeTeleBot.SimulateIncomingMessageToChat(1, "/test me")

	assert.Len(t, fakeTeleBot.OutboundSendMessages, 1)
	assert.Equal(t, `I matched with \/test (?P<who>\w+)`, fakeTeleBot.OutboundSendMessages[0])
}

func NewWrapBot(tBot tbwrap.TeleBot) *tbwrap.WrapBot {
	bot, err := tbwrap.NewBot(tbwrap.Config{
		TBot: tBot,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return bot
}
