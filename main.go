package main

import (
	"fmt"

	"github.com/jmmccray/weather-service/app"
)

func main() {
	fmt.Println("Starting OpenWeather service server...")
	app.StartOpenWeather()
}
