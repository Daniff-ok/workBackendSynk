package gateway

import (
	"net/http"
	"time"
	"os"
	"encoding/json"
	"fmt"
	"bytes"
	"io"

	"rndmBackend/internal/domain"
)


type TranslateGateway struct {
	url   string
	token string
	client http.Client
}

func NewTranslateGateway() *TranslateGateway {
	return &TranslateGateway{
		url: "https://api.deepseek.com/chat/completions",
		token:os.Getenv("DEEPSEEK_API"),
		client: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *TranslateGateway) RequestTranslate(text []string, language string) (domain.Response, error) {
	input, err := json.Marshal(text)
	if err != nil {
		return domain.Response{}, err
	}

	reqBody := domain.Request{
		Model: "deepseek-flash",
		Messages: []domain.Message{
			{
				Role:    "system",
				Content: fmt.Sprintf(`Translate each Dota 2 in-game player chat message to %s naturally, preserving gaming slang and abbreviations. Return minified JSON: {"translations":["..."]}. Preserve order..`, language),
			},
			{
				Role:    "user",
				Content: string(input),
			},
		},
	}

	reqBody.Thinking.Type = "disabled"
	reqBody.ResponseFormat.Type = "json_object"

	data, err := json.Marshal(reqBody)
	if err != nil {
		return domain.Response{}, err
	}

	req, err := http.NewRequest(
		"POST",
		r.url,
		bytes.NewReader(data),
	)
	if err != nil {
		return domain.Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

	resp, err := r.client.Do(req)
	if err != nil {
		return domain.Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return domain.Response{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return domain.Response{}, fmt.Errorf("DeepSeek API (%d): %s",
			resp.StatusCode, string(body))
	}

	var result domain.Response

	if err := json.Unmarshal(body, &result); err != nil {
		return domain.Response{}, err
	}

	return result, nil
}
