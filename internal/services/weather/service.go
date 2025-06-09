package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// {
//   "location": {
//     "name": "São Paulo",
//     "region": "Sao Paulo",
//     "country": "Brasilien",
//     "lat": -23.5333,
//     "lon": -46.6167,
//     "tz_id": "America/Sao_Paulo",
//     "localtime_epoch": 1749498880,
//     "localtime": "2025-06-09 16:54"
//   },
//   "current": {
//     "last_updated_epoch": 1749498300,
//     "last_updated": "2025-06-09 16:45",
//     "temp_c": 17.1,
//     "temp_f": 62.8,
//     "is_day": 1,
//     "condition": {
//       "text": "Partly cloudy",
//       "icon": "//cdn.weatherapi.com/weather/64x64/day/116.png",
//       "code": 1003
//     },
//     "wind_mph": 5.4,
//     "wind_kph": 8.6,
//     "wind_degree": 148,
//     "wind_dir": "SSE",
//     "pressure_mb": 1017,
//     "pressure_in": 30.03,
//     "precip_mm": 0.21,
//     "precip_in": 0.01,
//     "humidity": 72,
//     "cloud": 75,
//     "feelslike_c": 17.1,
//     "feelslike_f": 62.8,
//     "windchill_c": 15.9,
//     "windchill_f": 60.6,
//     "heatindex_c": 15.9,
//     "heatindex_f": 60.6,
//     "dewpoint_c": 14.6,
//     "dewpoint_f": 58.2,
//     "vis_km": 7,
//     "vis_miles": 4,
//     "uv": 0.1,
//     "gust_mph": 6.3,
//     "gust_kph": 10.1
//   }
// }

type WeatherResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float32 `json:"lat"`
		Lon            float32 `json:"lon"`
		TzId           string  `json:"tz_id"`
		LocaltimeEpoch int64   `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int64   `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float32 `json:"temp_c"`
		TempF            float32 `json:"temp_f"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
	}
}

func GetWeatherData(city string, apiKey string) (WeatherResponse, error) {
	if city == "" {
		fmt.Println("Error: city is empty")
		return WeatherResponse{}, fmt.Errorf("city is empty")
	}
	if apiKey == "" {
		fmt.Println("Error: API key is empty")
		return WeatherResponse{}, fmt.Errorf("API key is empty")
	}

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, url.QueryEscape(city))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching weather data:", err)
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	// DEBUG only
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return WeatherResponse{}, err
	// }
	// fmt.Println("Response body:", string(body))
	// END DEBUG

	if resp.Body == nil {
		fmt.Println("Error: response body is nil")
		return WeatherResponse{}, fmt.Errorf("response body is nil")
	}

	weatherResponse := WeatherResponse{}
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		fmt.Println("Error decoding weather response:", err)
		return WeatherResponse{}, err
	}

	return weatherResponse, nil
}
