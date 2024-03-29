package main

import(
	"github.com/jmmccray/weather-service/server"
	"fmt"
)

func main() {
	fmt.Println("Starting Weather Service server...")
	server.StartOpenWeatherApp()
}