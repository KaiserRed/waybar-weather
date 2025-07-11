package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Name string `json:"name"`
}

func getWeatherIcon(weatherType string) string {
	icons := map[string]string{
		"Thunderstorm": "⛈️",
		"Drizzle":      "🌧️",
		"Rain":         "🌧️",
		"Snow":         "❄️",
		"Clear":        "☀️",
		"Clouds":       "☁️",
	}
	if icon, ok := icons[weatherType]; ok {
		return icon
	}
	return "0"
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	BASE_URL := os.Getenv("BASE_URL")
	API_KEY := os.Getenv("API_KEY")
	var city string
	city = "Moscow"
	Search_Url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", BASE_URL, city, API_KEY)
	response, err := http.Get(Search_Url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatal(response.Status)
	}
	weatherBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var weather WeatherResponse
	if err := json.Unmarshal(weatherBytes, &weather); err != nil {
		log.Fatal(err)
	}

	output := map[string]interface{}{
		"text": fmt.Sprintf("%s %.1f°C", getWeatherIcon(weather.Weather[0].Main), weather.Main.Temp),
		"tooltip": fmt.Sprintf("%s: %s\nТемпература: %.1f°C\nОщущается как: %.1f°C\nВлажность: %d%%",
			weather.Name,
			weather.Weather[0].Description,
			weather.Main.Temp,
			weather.Main.FeelsLike,
			weather.Main.Humidity),
		"alt":   weather.Weather[0].Main,
		"class": strings.ToLower(weather.Weather[0].Main),
	}

	jsonOutput, _ := json.Marshal(output)
	fmt.Println(string(jsonOutput))
}
