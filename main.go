package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/rs/zerolog/log"
)

var (
	discordWebhookUrl = os.Getenv("DISCORD_WEBHOOK_URL")
)

type WAQIResponse struct {
	Status string `json:"status"`
	Data   struct {
		AQI int `json:"aqi"`
	} `json:"data"`
}

func stringToFloat(s string) (float64, error) {
	vInt, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, err
	}
	return vInt, nil
}

func getAQICategory(aqi int) (string, string) {
	switch {
	case aqi <= 50:
		return "🟢", "Good"
	case aqi <= 100:
		return "🟡", "Moderate"
	case aqi <= 150:
		return "🟠", "Unhealthy for Sensitive Groups"
	case aqi <= 200:
		return "🔴", "Unhealthy"
	case aqi <= 300:
		return "🟣", "Very Unhealthy"
	default:
		return "🟤", "Hazardous"
	}
}

func checkAqi(latitude, longitude float64) (string, error) {
	apiKey := os.Getenv("WAQI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("WAQI_API_KEY environment variable not set")
	}

	url := fmt.Sprintf("https://api.waqi.info/feed/geo:%.6f;%.6f/?token=%s", latitude, longitude, apiKey)

	var waqiResp WAQIResponse
	err := requests.URL(url).ToJSON(&waqiResp).Fetch(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to call WAQI API: %w", err)
	}

	if waqiResp.Status != "ok" {
		return "", fmt.Errorf("WAQI API returned error status: %s", waqiResp.Status)
	}

	emoji, category := getAQICategory(waqiResp.Data.AQI)
	return fmt.Sprintf("%s AQI: %d (%s)", emoji, waqiResp.Data.AQI, category), nil
}

func parseCoordinates() (float64, float64, error) {
	latitude, err := stringToFloat(os.Getenv("TARGET_LATITUDE"))
	if err != nil {
		return 0, 0, fmt.Errorf("error converting TARGET_LATITUDE to float: %w", err)
	}

	longitude, err := stringToFloat(os.Getenv("TARGET_LONGITUDE"))
	if err != nil {
		return 0, 0, fmt.Errorf("error converting TARGET_LONGITUDE to float: %w", err)
	}

	return latitude, longitude, nil
}

func main() {
	// parse coordinates from environment
	latitude, longitude, err := parseCoordinates()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse coordinates")
	}
	log.Info().Msgf("Latitude: %v", latitude)
	log.Info().Msgf("Longitude: %v", longitude)

	// check water
	outputMessage, err := checkAqi(latitude, longitude)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to check aqi info")
	}

	// send notification
	err = notify(outputMessage)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to send notification")
	}
	log.Info().Msg("Notification sent successfully")
}
