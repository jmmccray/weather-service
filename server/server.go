package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmmccray/weather-service/config"
)

type OpenWeatherClient struct {
	client *http.Client
}

type OpenWeatherServer struct {
	server *http.Server
}

type Coordinates struct {
	Latitude  string `json:"lat,omitempty"`
	Longitude string `json:"lon,omitempty"`
}

type Geolocation struct {
	City      string // name of city
	State     string
	Latitude  float64 // 6 decimal places
	Longitude float64 // 6 decimal places
}

// Defines the struct where the response from OpenWeather endpoint is stored.
type OpenWeatherData struct {
	Coord      OW_Coordinates `json:"coord,omitempty"`
	Weather    []WeatherData  `json:"weather,omitempty"`
	Base       string         `json:"base,omitempty"`
	Main       MainData       `json:"main,omitempty"`
	Visibility int            `json:"visibility,omitempty"`
	Wind       WindData       `json:"wind,omitempty"`
	Rain       RainData       `json:"rain,omitempty"`
	Clouds     CloudData      `json:"clouds,omitempty"`
	Dt         int            `json:"dt,omitempty"`
	System     SysData        `json:"sys,omitempty"`
	Timezone   int            `json:"timezone,omitempty"`
	Id         int            `json:"id,omitempty"`
	Name       string         `json:"name,omitempty"`
	Code       int            `json:"cod,omitempty"`
}

type OW_Coordinates struct {
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lon,omitempty"`
}

type WeatherData struct {
	Id          int    `json:"id,omitempty"`
	Main        string `json:"main,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type MainData struct {
	Temp        float64 `json:"temp,omitempty"`
	FeelsLike   float64 `json:"feels_like,omitempty"`
	TempMin     float64 `json:"temp_min,omitempty"`
	TempMax     float64 `json:"temp_max,omitempty"`
	Pressure    int     `json:"pressure,omitempty"`
	Humidity    int     `json:"humidity,omitempty"`
	SeaLevel    int     `json:"seal_level,omitempty"`
	GroundLevel int     `json:"grnd_level,omitempty"`
}

type WindData struct {
	Speed   float64 `json:"speed,omitempty"`
	Degrees int     `json:"deg,omitempty"`
	Gust    float64 `json:"gust,omitempty"`
}

type RainData struct {
	Hr float64 `json:"1h,omitempty"`
}

type CloudData struct {
	All int `json:"all,omitempty"`
}

type SysData struct {
	Type    int    `json:"type,omitempty"`
	Id      int    `json:"id,omitempty"`
	Country string `json:"country,omitempty"`
	Sunrise int    `json:"sunrise,omitempty"` // convert to unix to date format
	Sunset  int    `json:"sunset,omitempty"`
}

// TODO: Add OpenWeatherStatusCodes.
var openWeatherStatusCode map[int]string = map[int]string{}

// TODO: Elaborate on the Geolocation structs.
var geolocations []Geolocation = []Geolocation{
	{City: "San Diego", State: "CA", Latitude: 32.715736, Longitude: -117.161087},
	{City: "San Francisco", State: "CA", Latitude: 37.828724, Longitude: -122.355537},
	{City: "Houston", State: "TX", Latitude: 29.749907, Longitude: -95.358421},
}

// Create an OpenWeather client
func (c *OpenWeatherClient) NewClient() {
	c.client = &http.Client{}
}

func (s *OpenWeatherServer) NewServer() {
	s.server = &http.Server{}
}

// Make an HTTP GET request from OpenWeather client
func (c *OpenWeatherClient) ClientWeatherRequest(lat, lon string) error {
	// Check if client is initialize
	if c == nil {
		fmt.Println("Client has not be initialized")
		return errors.New("Client has not be initialized")
	}

	body := []byte(fmt.Sprintf("{\"lat\":\"%s\",\"lon\":\"%s\"}", lat, lon))

	// Preapre a http request to be sent to the OpenWeather server
	req, err := http.NewRequest("GET", "http://localhost:8080", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	var resp *http.Response
	for i := 1; i < 6; i++ {
		resp, err = c.client.Do(req)
		if err == nil {
			fmt.Println("OpenWeatherClient sent request")
			break
		}

		fmt.Println("Error: OpenWeatherClient couldn't send request:", err.Error())
		if i == 5 {
			return fmt.Errorf("error: OpenWeatherClient attempted 5 requests, exiting program")
		}
		// Tries to make request again after 2 seconds
		time.Sleep(2 * time.Second)
	}
	defer resp.Body.Close()
	fmt.Println()
	return nil
}

func (s *OpenWeatherServer) RunServer() {
	// Load necessary environment variables to get OpenWeather data
	err := config.LoadServerConfig()
	if err != nil {
		fmt.Println("Unable to load server config file")
		return
	}

	//http.HandleFunc("/client",getWeatherDataHandler)
	s.server = &http.Server{
		Addr:         "localhost:8080", // added localhost to by-pass Windows Defender firewall
		Handler:      http.HandlerFunc(getWeatherDataHandler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("OpenWeather server is running on %s\n", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil {
		panic(err)
	}
	
}

// Spins up a standalone the Weather Server instance on localhost:8080
func RunOpenWeatherServer(domain, port string) {
	// Load necessary environment variables to get OpenWeather data
	err := config.LoadServerConfig()
	if err != nil {
		fmt.Println("Unable to load server config file")
		return
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s",domain,port), // added localhost to by-pass Windows Defender firewall
		Handler:      http.HandlerFunc(getWeatherDataHandler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Running OpenWeather server on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func getWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	var weatherData OpenWeatherData
	var coordinates Coordinates
	var API_key = config.ConfigGetString("OW_API_KEY")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		w.Write([]byte(fmt.Sprintf("Error reading response body: %v\n", err)))
		return
	}

	if len(body) == 0 {
		fmt.Printf("Error: Request body is empty")
		w.Write([]byte(`Error: Request body is empty`))
		return
	}

	// Unmarshal data into OpenWeatherData object
	err = json.Unmarshal(body, &coordinates)
	if err != nil {
		fmt.Printf("Error unmarshaling data into Coordinates struct: %s\n", err.Error())
		w.Write([]byte(fmt.Sprintf("Error unmarshaling data into Coordinates struct: %s\n", err.Error())))
		return
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s",
		coordinates.Latitude,
		coordinates.Longitude,
		API_key)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Error sending request: %s\n", err.Error())
		w.Write([]byte(fmt.Sprintf("Error sending request: %s\n",err.Error())))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Status Code: ", resp.StatusCode)
		fmt.Printf("Invalid Coordinates: Lat=%s or Lon=%s\n", coordinates.Latitude, coordinates.Longitude)
		w.Write([]byte(fmt.Sprintf("Status Code: %d",resp.StatusCode)))
		return
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	w.Write([]byte(prettyFormatWeatherData(body)))
	fmt.Println("Weather Data: \n", prettyFormatWeatherData(body))

	// Unmarshal data into OpenWeatherData object
	err = json.Unmarshal(body, &weatherData)

	if err != nil {
		fmt.Printf("Error unmarshaling data into OpenWeatherData object: %s\n", err.Error())
		return
	}
	// TODO: write struct to SQL DB
}

func RunOpenWeatherClient() {
	client := &http.Client{}

	body := []byte(fmt.Sprintf("{\"lat\":\"%s\",\"lon\":\"%s\"}", "0", "0"))

	// Preapre a http request to be sent to the OpenWeather server
	req, err := http.NewRequest("GET", "http://localhost:8080", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	var resp *http.Response
	for i := 1; i <= 5; i++ {
		resp, err = client.Do(req)
		if err == nil {
			fmt.Println("OpenWeatherClient sent request")
			break
		}
		fmt.Println("Error: OpenWeatherClient couldn't send request to server", err)

		// Tries to make request again after 2 seconds
		time.Sleep(2 * time.Second)
		i++
	}
	fmt.Println()
	defer resp.Body.Close()
}

// Converts float geolocation coordinates into a string format.
func convertGeoCoors(coor float64) string {
	return fmt.Sprintf("%.6f", coor)
}

func validateCheckLatLonValue(lat, lon float64) error {
	if lat < -90 && lat > 90 {
		fmt.Printf("The latitude value '%f' value is not in the range of -90 <= lat <= 90\n", lat)
		return errors.New("invalid latitude type")
	} else if lon < -180 && lon > 180 {
		fmt.Printf("The longitude value '%f' value is not in the range of -180 <= lon <= 180\n", lon)
		return errors.New("invalid longitude type")
	}
	return nil
}

func getCLIPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	client := &OpenWeatherClient{}
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
			if ok := validateCheckLatLonValue(lat, lon); ok != nil {
				fmt.Printf("Invalid coordinates: %s\n", err.Error())
				continue
			}

			fmt.Printf("Lat: %s | Lon: %s\n", convertGeoCoors(lat), convertGeoCoors(lon))

			// Make HTTP request
			err = client.ClientWeatherRequest(convertGeoCoors(lat), convertGeoCoors(lon))
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}
			// TODO: implement select
		} else if userInput == "sel" {
			fmt.Println("Select from Menus:")

			// Generate selection menu
			menu := ""
			for i, loc := range geolocations {
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

			if index < 0 && index > len(geolocations)-1 {
				fmt.Printf("The number entered does exist in range [0-%d]\n",len(geolocations)-1)
			}
			client.ClientWeatherRequest(convertGeoCoors(geolocations[index].Latitude), convertGeoCoors(geolocations[index].Longitude))
		} else {
			fmt.Printf("Invalid input: %s\n", userInput)
		}
	}
}

// TODO: Define webapp prompt
func getWebAppPrompt() {

}

func StartOpenWeatherApp() {
	scanner := bufio.NewScanner(os.Stdin)
	server := &OpenWeatherServer{}
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

// Pretty-Prints the JSON response data in a read-able format to the console.
func prettyFormatWeatherData(data []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data, "", " ")
	if err != nil {
		fmt.Println("Error indenting JSON payload", err.Error())
	}
	return prettyJSON.String()

	// newdata, err := json.MarshalIndent(data, "", " ")
	// if err != nil {
	// 	fmt.Println("Error indenting JSON payload", err.Error())
	// }
	// return string(newdata)
}

// TODO: Prints the values from the OpenWeatherData struct
func printOpenWeatherDate(data OpenWeatherData) {

}
