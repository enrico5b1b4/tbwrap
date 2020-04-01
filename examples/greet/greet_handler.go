package main

import (
	"fmt"

	"github.com/enrico5b1b4/tbwrap"
)

type GreetMessage struct {
	Name string `regexpGroup:"name"`
}

func HandleGreet() func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(GreetMessage)
		if err := c.Bind(message); err != nil {
			return err
		}

		_, err := c.Send(fmt.Sprintf("Hello %s!", message.Name))

		return err
	}
}
