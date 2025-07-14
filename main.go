package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

func waitForInternet() {
	timeout := time.Now().Add(30 * time.Second)
	for time.Now().Before(timeout) {
		_, err := http.Get("http://clients3.google.com/generate_204")
		if err == nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Println("Warning: Internet connection not available after 30 seconds")
}

func loadEnv(path string) map[string]string {
	// file, err := os.Open(path)
	file, err := os.Open("/home/kaiser/.config/waybar/modules/.env")
	if err != nil {
		log.Fatal("Error oppening .env", err)
	}
	defer file.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			env[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading %v", err)
	}
	return env
}

func main() {
	waitForInternet()
	env := loadEnv(".env")
	BASE_URL := "http://api.openweathermap.org/data/2.5/weather"
	API_KEY := env["API_KEY"]
	city := env["CITY"]
	Search_Url := fmt.Sprintf("%v?q=%s&appid=%s&units=metric", BASE_URL, city, API_KEY)
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
