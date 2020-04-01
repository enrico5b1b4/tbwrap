package main

import (
	"fmt"

	"github.com/enrico5b1b4/tbwrap"
)

type WeatherMessage struct {
	Location string `regexpGroup:"location"`
}

func HandleWeather(weatherService *WeatherService) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(WeatherMessage)
		if err := c.Bind(message); err != nil {
			return err
		}

		forecast, err := weatherService.GetForecast(message.Location)
		if err != nil {
			return err
		}

		_, err = c.Send(fmt.Sprintf("%s %0.2f℃ (feels like %0.2f℃)", forecast.Icon, forecast.Temp, forecast.TempFeelsLike))

		return err
	}
}
