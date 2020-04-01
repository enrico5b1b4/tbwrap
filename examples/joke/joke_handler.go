package main

import (
	"encoding/json"
	"net/http"

	"github.com/enrico5b1b4/tbwrap"
)

type Joke struct {
	ID   string `json:"id"`
	Joke string `json:"joke"`
}

func HandleJoke() func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		req, err := http.NewRequest(http.MethodGet, "https://icanhazdadjoke.com/", nil)
		if err != nil {
			return err
		}
		req.Header.Add("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var joke Joke
		err = json.NewDecoder(resp.Body).Decode(&joke)
		if err != nil {
			return err
		}

		_, err = c.Send(joke.Joke)

		return err
	}
}
