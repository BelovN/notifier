package meteo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BelovN/notifier/internal/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type CurrentWeather struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windspeed"`
	WeatherCode int     `json:"weathercode"`
}

func (cw CurrentWeather) ToString() string {
	description, exists := WeatherCodeDescriptions[cw.WeatherCode]
	var weatherDescription string
	if !exists {
		weatherDescription = ""
	} else {
		weatherDescription = fmt.Sprintf("Описание: %s", description)
	}
	return fmt.Sprintf(
		"Current weather: temperature: %.2f °C, wind speed: %.2f km/h. %s.",
		cw.Temperature,
		cw.WindSpeed,
		weatherDescription,
	)
}

type Response struct {
	CurrentWeather CurrentWeather `json:"current_weather"`
}

type ServiceConfig struct {
	Latitude  string
	Longitude string
	BaseUrl   string
	Timeout   time.Duration
}

type Service struct {
	Client *http.Client
	Ctx    context.Context

	Cfg ServiceConfig
}

func NewService(ctx context.Context, cfg *config.Config) *Service {
	return &Service{Client: &http.Client{}, Ctx: ctx, Cfg: ServiceConfig{
		cfg.Meteo.Latitude, cfg.Meteo.Longitude, cfg.Meteo.BaseUrl, cfg.Meteo.Timeout,
	}}
}

func (s *Service) makeRequest(params url.Values) ([]byte, error) {
	ctx, cancel := context.WithTimeout(s.Ctx, s.Cfg.Timeout)
	defer cancel()

	fullURL := s.Cfg.BaseUrl + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errors.New("failed to execute request: " + err.Error())
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	return body, nil
}

func (s *Service) GetCurrentWeather() (*CurrentWeather, error) {
	params := url.Values{
		"latitude":        {s.Cfg.Latitude},
		"longitude":       {s.Cfg.Longitude},
		"current_weather": {"true"},
		"windspeed_unit":  {"kmh"},
	}

	body, err := s.makeRequest(params)
	if err != nil {
		return nil, err
	}

	var weatherResponse Response
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return &weatherResponse.CurrentWeather, nil
}
