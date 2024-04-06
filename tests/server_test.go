package tests

import (
	"net/http"
	"net/http/httptest"
	s "github.com/jmmccray/weather-service/server"
	"testing"
)

func TestClientWeatherRequest(t *testing.T) {
    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Respond with a mock response
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"key": "value"}`))
    }))
    defer server.Close()

    // Create a new OpenWeatherClient
    client := &s.OpenWeatherClient{}
    client.NewClient()

    // Make a request to the mock server
    err := client.ClientWeatherRequest("0", "0", server.URL)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
}