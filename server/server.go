package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"bytes"
	"github.com/jmmccray/weather-service/config"
)

type Geolocation struct {
	City      string // name of city
	State     string
	Latitude  float64 // 6 decimal places
	Longitude float64 // 6 decimal places
}

// Defines the struct where the response from OpenWeather endpoint is stored.
type OpenWeatherData struct {
	Coord      Coordinate    `json:"coord,omitempty"`
	Weather    []WeatherData `json:"weather,omitempty"`
	Base       string        `json:"base,omitempty"`
	Main       MainData      `json:"main,omitempty"`
	Visibility int        `json:"visibility,omitempty"`
	Wind       WindData      `json:"wind,omitempty"`
	Rain       RainData      `json:"rain,omitempty"`
	Clouds     CloudData     `json:"clouds,omitempty"`
	Dt         int        `json:"dt,omitempty"`
	System     SysData       `json:"sys,omitempty"`
	Timezone   int           `json:"timezone,omitempty"`
	Id         int           `json:"id,omitempty"`
	Name       string        `json:"name,omitempty"`
	Code       int           `json:"cod,omitempty"`
}

type Coordinate struct {
	Longitude float64 `json:"lon,omitempty"`
	Latitude  float64 `json:"lat,omitempty"`
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
var OpenWeatherStatusCode map[int]string = map[int]string{}

// TODO: Elaborate on the Geolocation structs.
var Geolocations []Geolocation = []Geolocation{
	{City: "San Diego", State: "CA", Latitude: 32.715736, Longitude: -117.161087},
	{City: "San Francisco", State: "CA", Latitude: 37.828724, Longitude: -122.355537},
	{City: "Houston", State: "TX", Latitude: 29.749907, Longitude: -95.358421}}

var API_key = "41ed40847bc3fa78c22ac0f73bef1477"

// Converts geolocation coordinates into the coordinate defined string format.
func convertGeoCoors(coor float64) string {
	return fmt.Sprintf("%.6f", coor)
}

// Spins up the Weather Server on localhost:8080
func RunWeatherServer() {
	// Load necessary environment variables to get OpenWeather data
	err := config.LoadServerConfig()
	if err != nil {
		fmt.Println("Unable to load server config file")
		return
	}

	server := &http.Server{
		Addr:         "localhost:8080", // added localhost to by-pass Windows Defender firewall
		Handler:      http.HandlerFunc(getWeatherDataHandler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Running on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func OpenWeatherClient() {
	url := "http://localhost:8080"

	for {
		req, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}
		_ = req
	}
	
	// TODO: set headers
}

func getWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Calling getWeatherDataHandler")
	var weatherData OpenWeatherData
	var API_key = config.ConfigGetString("OW_API_KEY")

	fmt.Println(API_key)

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s",
		convertGeoCoors(Geolocations[0].Latitude),
		convertGeoCoors(Geolocations[0].Longitude),
		API_key)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Error sending request: %v\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Status Code: ", resp.StatusCode)
		fmt.Println("Invalid Coordinates: Lat=%d or Lon=%d", 0, 0)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	w.Write([]byte(prettyFormatWeatherData([]byte(body))))

	fmt.Println("Weather Data: \n", prettyFormatWeatherData([]byte(body)))

	// Marshal data into OpenWeatherData object
	err  = json.Unmarshal(body, &weatherData)

	if err != nil {
		fmt.Printf("Error unmarshaling data into OpenWeatherData object: %s\n", err.Error())
	}

	// TODO: write struct to SQL DB
}

// Pretty-Prints the JSON response data in a read-able format to the console.
func prettyFormatWeatherData(data []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data,"", " ")
	if err != nil {
		fmt.Println("Error identing JSON payload", err.Error())
	}
	return prettyJSON.String()
}

// TODO: Prints the values from the OpenWeatherData struct
func printOpenWeatherDate(data OpenWeatherData) {
	 
}
