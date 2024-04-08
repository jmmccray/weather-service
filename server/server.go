package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jmmccray/weather-service/config"
	"github.com/jmmccray/weather-service/models"
	"github.com/jmmccray/weather-service/utils"
)

type OpenWeatherServer struct {
	server *http.Server
}

func (s *OpenWeatherServer) NewServer() {
	s.server = &http.Server{}
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
	var weatherData models.OpenWeatherData
	var coordinates models.Coordinates
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

	w.Write([]byte(utils.PrettyFormatWeatherData(body)))
	fmt.Println("Weather Data: \n", utils.PrettyFormatWeatherData(body))

	err = json.Unmarshal(body, &weatherData)

	if err != nil {
		fmt.Printf("Error unmarshaling data into OpenWeatherData object: %s\n", err.Error())
		return
	}
	// TODO: write struct to SQL DB
}
