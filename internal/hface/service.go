package hface

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type HfaceService struct {
	ctx    context.Context
	client *http.Client
	token  string
}

const (
	BaseURL = "https://api-inference.huggingface.co/models/meta-llama/Meta-Llama-3-8B-Instruct/v1/chat/completions"
	Timeout = 10 * time.Second
)

func NewHfaceService(ctx context.Context, token string, client *http.Client) *HfaceService {
	if client == nil {
		client = &http.Client{}
	}
	return &HfaceService{ctx: ctx, token: token, client: client}
}

func (s *HfaceService) GetAIAnswer(content string) (string, error) {
	randomNumber := rand.Intn(1000)

	messages := []HfaceMessage{
		SystemHfaceMessage,
		{Role: "user", Content: content},
	}
	payload := HfaceRequestPayload{messages, 500, false, randomNumber}

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
