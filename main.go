package main

import (
	"fmt"

	"github.com/jmmccray/weather-service/server"
)

func main() {
	fmt.Println("Starting OpenWeather service server...")
	server.StartOpenWeatherApp()
}
