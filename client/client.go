package client

import (
	"net/http"
	"fmt"
	"bytes"
	"errors"
	"time"
)

type OpenWeatherClient struct {
	client *http.Client
}

// Create an OpenWeather client
func (c *OpenWeatherClient) NewClient() {
	c.client = &http.Client{}
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