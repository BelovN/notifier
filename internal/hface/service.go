package hface

import (
	"context"
	"encoding/json"
	"github.com/BelovN/notifier/internal/config"
	"math/rand"
	"net/http"
	"time"
)

type ServiceConfig struct {
	BaseUrl string
	Timeout time.Duration
}

type Service struct {
	Ctx    context.Context
	Client *http.Client
	Token  string

	Cfg ServiceConfig
}

func NewService(ctx context.Context, token string, cfg *config.Config) *Service {
	return &Service{ctx, &http.Client{}, token,
		ServiceConfig{cfg.HFace.BaseUrl, cfg.HFace.Timeout},
	}
}

func (s *Service) GetAIAnswer(content string) (string, error) {
	randomNumber := rand.Intn(1000)

	messages := []Message{
		SystemHfaceMessage,
		{Role: "user", Content: content},
	}
	payload := RequestPayload{messages, 500, false, randomNumber}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	data, err := s.makeRequest(jsonData)
	if err != nil {
		return "", err
	}

	return data.Choices[0].Message.Content, nil
}
