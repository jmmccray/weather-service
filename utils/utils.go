package utils

import (
	"fmt"
	"errors"
	"bytes"
	"encoding/json"
	"github.com/jmmccray/weather-service/models"
)

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

func PrettyFormatWeatherData(data []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data, "", " ")
	if err != nil {
		fmt.Println("Error indenting JSON payload", err.Error())
	}
	return prettyJSON.String()
}

func PrintOpenWeatherData(data models.OpenWeatherData) {

}
