package tbwrap

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func NewPollerWithAllowedChats(pollTimout time.Duration, chats []int) *tb.MiddlewarePoller {
	poller := &tb.LongPoller{Timeout: pollTimout}

	return tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		// allow request from any user/group if no restrictions are set
		if len(chats) == 0 {
			return true
		}

		if upd.Message != nil {
			return isInList(int(upd.Message.Chat.ID), chats)
		}

		if upd.Callback != nil {
			return isInList(int(upd.Callback.Message.Chat.ID), chats)
		}

		return false
	})
}

func isInList(id int, list []int) bool {
	for i := range list {
		if list[i] == id {
			return true
		}
	}
	return false
}
