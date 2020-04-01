package main

import (
	"bytes"
	"html/template"

	"github.com/enrico5b1b4/tbwrap"
)

func HandleList(todos map[int64][]string) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		chatTodos := todos[c.ChatID()]
		if len(chatTodos) == 0 {
			_, err := c.Send("your todo list is empty")

			return err
		}

		t := template.Must(template.New("text").Parse(text))
		var buf bytes.Buffer
		if err := t.Execute(&buf, chatTodos); err != nil {
			return err
		}

		_, err := c.Send(buf.String())

		return err
	}
}

const text = `{{ range $i, $entry := . }}{{printf "%d - %s\n" $i $entry}}{{ end }}`
