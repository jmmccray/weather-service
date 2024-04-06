package app

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"time"
	"strconv"

	"github.com/jmmccray/weather-service/server"
	"github.com/jmmccray/weather-service/client"
	"github.com/jmmccray/weather-service/models"
)

func StartOpenWeather() {
	scanner := bufio.NewScanner(os.Stdin)
	server := &server.OpenWeatherServer{}
	go server.RunServer()
	time.Sleep(3 * time.Second)

	// Allows the user to input an incorrect response 5x.
	for i := 1; i <= 5; {
		// Check if the user wants to use the cli or web app interfaces.
		fmt.Print("Would you like to use the CLI or web app (cli/app)? ") // Fix wording here
		scanner.Scan()

		// Gets user input, removes all of its whitespace and lowercases each character.
		userInput := strings.ToLower(strings.ReplaceAll(scanner.Text(), " ", ""))

		if userInput == "cli" {
			getCLIPrompt()
		} else if userInput == "app" {
			getWebAppPrompt()
		} else {
			fmt.Printf("Invalid input: %s\n", userInput)
			i++
			if i > 5 {
				fmt.Println("Too many invalid attempts, exiting program.")
				os.Exit(0)
			}
		}
	}
}

func getCLIPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	client := &client.OpenWeatherClient{}
	client.NewClient()

	for {
		// Ask if the user wants to enter coordinates manually or choose from preset list.
		fmt.Println("Manually enter latitude and longitude coordinates (man) or") // Fix wording here
		fmt.Print("Select from a predefined list (sel)? ")
		scanner.Scan()

		userInput := strings.ToLower(strings.ReplaceAll(scanner.Text(), " ", ""))

		if userInput == "man" {
			fmt.Println("Enter your coorindates with format:")
			fmt.Println("1. Format: XX.XXXXXX")
			fmt.Println("2. Latitude: -90 <= latitude <= 90 | Longitude: -180 <= longitude <= 180")

			var lat, lon float64
			var err error

			// Take user input for latitude value
			fmt.Print("Latitude: ")
			scanner.Scan()
			latitude := scanner.Text()
			// Check if latitude string can be converted to a float, then convert to float
			if lat, err = strconv.ParseFloat(latitude, 64); err != nil {
				fmt.Println("ERROR: You inputted in an invalid latitude value: ", latitude)
				continue
			}

			// Take user input for longitude value
			fmt.Print("Longitude: ")
			scanner.Scan()
			longitude := scanner.Text()
			// Check if longitude can be converted to a float, then convert to float
			if lon, err = strconv.ParseFloat(latitude, 64); err != nil {
				fmt.Println("ERROR: You inputted in an invalid longitude value: ", longitude)
				continue
			}

			// Check if lat and lon are valid coordinates
			if ok := models.ValidateCheckLatLonValue(lat, lon); ok != nil {
				fmt.Printf("Invalid coordinates: %s\n", err.Error())
				continue
			}

			fmt.Printf("Lat: %s | Lon: %s\n", models.ConvertGeoCoors(lat), models.ConvertGeoCoors(lon))

			// Make HTTP request
			err = client.ClientWeatherRequest(models.ConvertGeoCoors(lat), models.ConvertGeoCoors(lon))
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}
			// TODO: implement select
		} else if userInput == "sel" {
			fmt.Println("Select from Menus:")

			// Generate selection menu
			menu := ""
			for i, loc := range models.Geolocations {
				menu += fmt.Sprintf("%d. %s, %s\n",i, loc.City,loc.State)
			}
			fmt.Println(menu)
			
			fmt.Print("Enter #: ")
			scanner.Scan()

			// Gets user input, removes all of its whitespace and lowercases each character.
			userInput := scanner.Text()

			index, err := strconv.Atoi(userInput)
			if err != nil {
				fmt.Println("The input is not a valid number: ", index)
			}

			if index < 0 && index > len(models.Geolocations)-1 {
				fmt.Printf("The number entered does exist in range [0-%d]\n",len(models.Geolocations)-1)
			}
			client.ClientWeatherRequest(models.ConvertGeoCoors(models.Geolocations[index].Latitude), models.ConvertGeoCoors(models.Geolocations[index].Longitude))
		} else {
			fmt.Printf("Invalid input: %s\n", userInput)
		}
	}
}

// TODO: Define webapp prompt
func getWebAppPrompt() {

}


