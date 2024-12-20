package hface

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestPayload struct {
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream"`
	Seed      int       `json:"seed"`
}

var SystemHfaceMessage = Message{
	Role: "system",
	Content: `You are my assistant. Say hello. Provide me with a concise summary of the current weather, 
including temperature, precipitation, and other relevant details. 
Add a humorous comment or joke related to the weather, and finish with a practical piece of advice 
to make the most of the day based on the conditions.`,
}

type ResponseMessage struct {
	Content string `json:"content"`
}

type ResponseChoice struct {
	Message ResponseMessage `json:"message"`
}

type Response struct {
	Choices []ResponseChoice `json:"choices"`
}

func (s *Service) makeRequest(jsonData []byte) (*Response, error) {
	ctx, cancel := context.WithTimeout(s.Ctx, s.Cfg.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", s.Cfg.BaseUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errors.New("failed to execute request: " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unknown status: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	var hfaceResponse Response
	if err := json.Unmarshal(body, &hfaceResponse); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return &hfaceResponse, nil
}
