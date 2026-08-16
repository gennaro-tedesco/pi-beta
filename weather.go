package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const ipLocationURL = "http://ip-api.com/json/"
const weatherForecastURL = "https://api.open-meteo.com/v1/forecast"
const airQualityURL = "https://air-quality-api.open-meteo.com/v1/air-quality"
const reverseGeocodeURL = "https://nominatim.openstreetmap.org/reverse"
const reverseGeocodeUserAgent = "pi-beta-weather-dashboard"
const reverseGeocodeZoom = 10
const weatherHTTPTimeout = 5 * time.Second
const invalidUVDataMessage = "forecast response does not contain UV data for the current time"
const currentDayIndex = 0
const hourlyForecastWindow = 24
const weatherForecastQueryFormat = "%s?latitude=%f&longitude=%f&current=temperature_2m,weather_code,is_day,wind_speed_10m,precipitation_probability,surface_pressure&hourly=uv_index,temperature_2m,weather_code&daily=weather_code,temperature_2m_max,temperature_2m_min,sunrise,sunset&forecast_days=8&timezone=auto"
const airQualityQueryFormat = "%s?latitude=%f&longitude=%f&current=us_aqi"

var weatherHTTPClient = &http.Client{Timeout: weatherHTTPTimeout}

type Weather struct {
	City                string    `json:"city"`
	Temperature         float64   `json:"temperature"`
	WeatherCode         int       `json:"weatherCode"`
	IsDay               bool      `json:"isDay"`
	WindSpeed           float64   `json:"windSpeed"`
	RainProbability     int       `json:"rainProbability"`
	UVIndex             float64   `json:"uvIndex"`
	Pressure            float64   `json:"pressure"`
	AirQuality          float64   `json:"airQuality"`
	DailyTime           []string  `json:"dailyTime"`
	DailyWeatherCode    []int     `json:"dailyWeatherCode"`
	DailyTemperatureMax []float64 `json:"dailyTemperatureMax"`
	DailyTemperatureMin []float64 `json:"dailyTemperatureMin"`
	Sunrise             string    `json:"sunrise"`
	Sunset              string    `json:"sunset"`
	HourlyTime          []string  `json:"hourlyTime"`
	HourlyTemperature   []float64 `json:"hourlyTemperature"`
	HourlyWeatherCode   []int     `json:"hourlyWeatherCode"`
	UTCOffsetSeconds    int       `json:"utcOffsetSeconds"`
}

type ipLocation struct {
	City string  `json:"city"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type nominatimResponse struct {
	DisplayName string `json:"display_name"`
	Address     struct {
		City    string `json:"city"`
		Town    string `json:"town"`
		Village string `json:"village"`
	} `json:"address"`
}

type openMeteoResponse struct {
	UTCOffsetSeconds int `json:"utc_offset_seconds"`
	Current          struct {
		Time            string  `json:"time"`
		Temperature2m   float64 `json:"temperature_2m"`
		WeatherCode     int     `json:"weather_code"`
		IsDay           int     `json:"is_day"`
		WindSpeed10m    float64 `json:"wind_speed_10m"`
		RainProbability int     `json:"precipitation_probability"`
		SurfacePressure float64 `json:"surface_pressure"`
	} `json:"current"`
	Hourly struct {
		Time          []string  `json:"time"`
		UVIndex       []float64 `json:"uv_index"`
		Temperature2m []float64 `json:"temperature_2m"`
		WeatherCode   []int     `json:"weather_code"`
	} `json:"hourly"`
	Daily struct {
		Time             []string  `json:"time"`
		WeatherCode      []int     `json:"weather_code"`
		Temperature2mMax []float64 `json:"temperature_2m_max"`
		Temperature2mMin []float64 `json:"temperature_2m_min"`
		Sunrise          []string  `json:"sunrise"`
		Sunset           []string  `json:"sunset"`
	} `json:"daily"`
}

type airQualityResponse struct {
	Current struct {
		USAQI float64 `json:"us_aqi"`
	} `json:"current"`
}

func (a *App) GetWeather() (Weather, error) {
	location, err := fetchIPLocation()
	if err != nil {
		return Weather{}, err
	}

	return buildWeather(location.City, location.Lat, location.Lon)
}

func (a *App) GetWeatherAt(lat, lon float64) (Weather, error) {
	city, err := reverseGeocode(lat, lon)
	if err != nil {
		return Weather{}, err
	}

	return buildWeather(city, lat, lon)
}

func buildWeather(city string, lat, lon float64) (Weather, error) {
	forecast, err := fetchForecast(lat, lon)
	if err != nil {
		return Weather{}, err
	}
	airQuality, err := fetchAirQuality(lat, lon)
	if err != nil {
		return Weather{}, err
	}

	result := Weather{
		City:                city,
		Temperature:         forecast.Current.Temperature2m,
		WeatherCode:         forecast.Current.WeatherCode,
		IsDay:               forecast.Current.IsDay == 1,
		WindSpeed:           forecast.Current.WindSpeed10m,
		RainProbability:     forecast.Current.RainProbability,
		Pressure:            forecast.Current.SurfacePressure,
		AirQuality:          airQuality.Current.USAQI,
		DailyTime:           forecast.Daily.Time,
		DailyWeatherCode:    forecast.Daily.WeatherCode,
		DailyTemperatureMax: forecast.Daily.Temperature2mMax,
		DailyTemperatureMin: forecast.Daily.Temperature2mMin,
		Sunrise:             forecast.Daily.Sunrise[currentDayIndex],
		Sunset:              forecast.Daily.Sunset[currentDayIndex],
		UTCOffsetSeconds:    forecast.UTCOffsetSeconds,
	}

	currentHour := forecast.Current.Time
	if len(currentHour) >= 13 {
		currentHour = currentHour[:13] + ":00"
	}
	currentIndex := indexOf(forecast.Hourly.Time, currentHour)
	if currentIndex < 0 || currentIndex >= len(forecast.Hourly.UVIndex) {
		return Weather{}, errors.New(invalidUVDataMessage)
	}
	result.UVIndex = forecast.Hourly.UVIndex[currentIndex]

	endIndex := currentIndex + hourlyForecastWindow
	if endIndex > len(forecast.Hourly.Time) {
		endIndex = len(forecast.Hourly.Time)
	}
	result.HourlyTime = forecast.Hourly.Time[currentIndex:endIndex]
	result.HourlyTemperature = forecast.Hourly.Temperature2m[currentIndex:endIndex]
	result.HourlyWeatherCode = forecast.Hourly.WeatherCode[currentIndex:endIndex]

	return result, nil
}

func reverseGeocode(lat, lon float64) (string, error) {
	url := fmt.Sprintf("%s?format=json&lat=%f&lon=%f&zoom=%d", reverseGeocodeURL, lat, lon, reverseGeocodeZoom)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", reverseGeocodeUserAgent)

	resp, err := weatherHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("reverse geocode request failed: %s", resp.Status)
	}

	var result nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Address.City != "" {
		return result.Address.City, nil
	}
	if result.Address.Town != "" {
		return result.Address.Town, nil
	}
	if result.Address.Village != "" {
		return result.Address.Village, nil
	}

	return result.DisplayName, nil
}

func fetchIPLocation() (ipLocation, error) {
	resp, err := weatherHTTPClient.Get(ipLocationURL)
	if err != nil {
		return ipLocation{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ipLocation{}, fmt.Errorf("ip location request failed: %s", resp.Status)
	}

	var location ipLocation
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return ipLocation{}, err
	}

	return location, nil
}

func fetchForecast(lat, lon float64) (openMeteoResponse, error) {
	url := fmt.Sprintf(
		weatherForecastQueryFormat,
		weatherForecastURL, lat, lon,
	)

	resp, err := weatherHTTPClient.Get(url)
	if err != nil {
		return openMeteoResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return openMeteoResponse{}, fmt.Errorf("forecast request failed: %s", resp.Status)
	}

	var forecast openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return openMeteoResponse{}, err
	}

	return forecast, nil
}

func fetchAirQuality(lat, lon float64) (airQualityResponse, error) {
	url := fmt.Sprintf(airQualityQueryFormat, airQualityURL, lat, lon)
	resp, err := weatherHTTPClient.Get(url)
	if err != nil {
		return airQualityResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return airQualityResponse{}, fmt.Errorf("air quality request failed: %s", resp.Status)
	}

	var airQuality airQualityResponse
	if err := json.NewDecoder(resp.Body).Decode(&airQuality); err != nil {
		return airQualityResponse{}, err
	}

	return airQuality, nil
}

func indexOf(values []string, target string) int {
	for i, value := range values {
		if value == target {
			return i
		}
	}
	return -1
}
