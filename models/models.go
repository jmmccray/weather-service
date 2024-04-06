package models

import (
	"encoding/json"
	"bytes"
	"fmt"

	"errors"
)

// TODO: Add OpenWeatherStatusCodes.
var openWeatherStatusCode map[int]string = map[int]string{}

// TODO: Elaborate on the Geolocation structs.
var Geolocations []Geolocation = []Geolocation{
	{City: "San Diego", State: "CA", Latitude: 32.715736, Longitude: -117.161087},
	{City: "San Francisco", State: "CA", Latitude: 37.828724, Longitude: -122.355537},
	{City: "Houston", State: "TX", Latitude: 29.749907, Longitude: -95.358421},
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

// Converts float geolocation coordinates into a string format.
func ConvertGeoCoors(coor float64) string {
	return fmt.Sprintf("%.6f", coor)
}

func ValidateCheckLatLonValue(lat, lon float64) error {
	if lat < -90 && lat > 90 {
		fmt.Printf("The latitude value '%f' value is not in the range of -90 <= lat <= 90\n", lat)
		return errors.New("invalid latitude type")
	} else if lon < -180 && lon > 180 {
		fmt.Printf("The longitude value '%f' value is not in the range of -180 <= lon <= 180\n", lon)
		return errors.New("invalid longitude type")
	}
	return nil
}

// Pretty-Prints the JSON response data in a read-able format to the console.
func PrettyFormatWeatherData(data []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data, "", " ")
	if err != nil {
		fmt.Println("Error indenting JSON payload", err.Error())
	}
	return prettyJSON.String()
}

// TODO: Prints the values from the OpenWeatherData struct
func PrintOpenWeatherData(data OpenWeatherData) {

}