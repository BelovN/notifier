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

type HfaceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type HfaceRequestPayload struct {
	Messages  []HfaceMessage `json:"messages"`
	MaxTokens int            `json:"max_tokens"`
	Stream    bool           `json:"stream"`
	Seed      int            `json:"seed"`
}

var SystemHfaceMessage = HfaceMessage{
	Role:    "system",
	Content: "You are my assistant. Say hello. Provide me with a concise summary of the current weather, including temperature, precipitation, and other relevant details. Add a humorous comment or joke related to the weather, and finish with a practical piece of advice to make the most of the day based on the conditions.",
}

type HfaceResponseMessage struct {
	Content string `json:"content"`
}

type HfaceResponseChoice struct {
	Message HfaceResponseMessage `json:"message"`
}

type HfaceResponse struct {
	Choices []HfaceResponseChoice `json:"choices"`
}

func (s *HfaceService) makeRequest(jsonData []byte) (*HfaceResponse, error) {
	ctx, cancel := context.WithTimeout(s.ctx, Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
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

	var hfaceResponse HfaceResponse
	if err := json.Unmarshal(body, &hfaceResponse); err != nil {
		return nil, errors.New("failed to decode response: " + err.Error())
	}

	return &hfaceResponse, nil
}
