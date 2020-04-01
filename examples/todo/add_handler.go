package main

import (
	"fmt"

	"github.com/enrico5b1b4/tbwrap"
)

type AddMessage struct {
	Value string `regexpGroup:"value"`
}

func HandleAdd(todos map[int64][]string) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(AddMessage)
		if err := c.Bind(message); err != nil {
			return err
		}

		todos[c.ChatID()] = append(todos[c.ChatID()], message.Value)

		_, err := c.Send(fmt.Sprintf(`"%s" has been added to your todo list`, message.Value))

		return err
	}
}
