# tbwrap
Wrapper library for [github.com/tucnak/telebot](https://github.com/tucnak/telebot) to create simple text-based Telegram bots

## Installation

```
go get github.com/enrico5b1b4/tbwrap
```

## Examples

[Ping Bot](./examples/ping) - A bot which responds to your ping with a pong  
[Greet Bot](./examples/greet) - A bot which greets you  
[Todo List Bot](./examples/todo) - A bot that keeps a todo list for you  
[Joke Bot](./examples/joke) - A bot that tells you a joke  
[Weather Bot](./examples/weather) - A bot that tells you the weather forecast  

## Usage
### Simple
```go
package main

import (
	"log"

	"github.com/enrico5b1b4/tbwrap"
)

func main() {
	telegramBotToken := "TELEGRAM_BOT_TOKEN"

	botConfig := tbwrap.Config{
		Token: telegramBotToken,
	}
	telegramBot, err := tbwrap.NewBot(botConfig)
	if err != nil {
		log.Fatal(err)
	}

	telegramBot.Handle(`/ping`, func(c tbwrap.Context) error {
		_, err := c.Send("pong!")

		return err
	})
	telegramBot.Start()
}
```
### Handlers
There are three ways you can register handlers to respond to a user's message

### Handle
Invokes handler if message matches the text provided
```go
// ...
func main() {
	// ...
	telegramBot.Handle(`/ping`, func(c tbwrap.Context) error {
		// ...
		_, err := c.Send("...")

		return err
	})
	// ...
}
```

### HandleRegExp
Invokes handler if message matches the regular expression provided
```go
// ...
func main() {
	// ...
	telegramBot.HandleRegExp(`\/greet (?P<name>.*)`, func(c tbwrap.Context) error {
		// ...
		_, err := c.Send("...")

		return err
    })
	// ...
}
```
### HandleMultiRegExp
Invokes handler if message matches any of the regular expressions provided
```go
// ...
func main() {
	// ...
	telegramBot.HandleMultiRegExp([]string{
		`\/greet (?P<name>.*)`,
		`\/welcome (?P<name>.*)`,
	}, func(c tbwrap.Context) error {
		// ...
		_, err := c.Send("...")

		return err
    })
	// ...
}
```
### Messages  
Read the user's incoming message with `c.Text()`  
```go
// ...

func main() {
	// ...
	telegramBot.HandleRegExp(`\/greet (?P<name>.*)`, func(c tbwrap.Context) error {
		// ...
		fmt.Println(c.Text()) // "/greet Enrico"

		_, err := c.Send("Hello!")

		return err
    })
	// ...
}
```
Use regular expression named capturing groups to parametrise incoming messages
```go
// ...

type GreetMessage struct {
    Name string `regexpGroup:"name"`
}

func main() {
	// ...
	telegramBot.HandleRegExp(`\/greet (?P<name>.*)`, func(c tbwrap.Context) error {
		message := new(GreetMessage)
		if err := c.Bind(message); err != nil {
			return err
		}

		fmt.Println(message.Name) // "Enrico"

		_, err := c.Send(fmt.Sprintf("Hello %s!", message.Name))

		return err
    })
	// ...
}
```
### Private bot
You can make the bot respond only to certain chats (users or groups) by passing a list of ids.  

```go
// ...
func main() {
	// ...
	botConfig := tbwrap.Config{
		Token:  telegramBotToken,
		AllowedChats:  []int{"123456","234567","-999999"}, // only user and group chats with these ids can interact with the bot
	}
	telegramBot, err := tbwrap.NewBot(botConfig)
	if err != nil {
		log.Fatal(err)
	}
	// ...
}
```