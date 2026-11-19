package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Thinking struct {
		Type string `json:"type"`
	} `json:"thinking"`
	ResponseFormat struct {
		Type string `json:"type"`
	} `json:"response_format"`
}

type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

var client = &http.Client{
	Timeout: 30 * time.Second,
}

func returnTranslateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	var request TRrequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65_536)).Decode(&request); err != nil ||
		len(request.Items) == 0 || len(request.Items) > 80 {
		http.Error(w, "Expected 1 to 80 messages", http.StatusBadRequest)
		return
	}
	language := strings.TrimSpace(request.TargetLang)
	if language == "" {
		language = "RU"
	} else if len(language) > 16 || strings.ContainsFunc(language, func(ch rune) bool {
		return (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') && ch != '-'
	}) {
		http.Error(w, "Invalid target language", http.StatusBadRequest)
		return
	}
	text := make([]string, len(request.Items))
	for i, item := range request.Items {
		text[i] = item.Text
	}
	translated, err := Translate(r.Context(), text, language)
	if err != nil {
		http.Error(w, "Translation provider unavailable", http.StatusBadGateway)
		return
	}
	response := TranslateReturn{Results: make([]TranslateText, len(translated))}
	for i, value := range translated {
		response.Results[i] = TranslateText{value, request.Items[i].SourceLanguage}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(response)
}

func Translate(ctx context.Context, text []string, language string) ([]string, error) {
	if len(text) == 0 {
		return []string{}, nil
	}

	input, err := json.Marshal(text)
	if err != nil {
		return nil, err
	}

	reqBody := Request{
		Model: "deepseek-flash",
		Messages: []Message{
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
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.deepseek.com/chat/completions",
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("DEEPSEEK_API"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DeepSeek API (%d): %s",
			resp.StatusCode, string(body))
	}

	var result Response

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	var translated struct {
		Translations []string `json:"translations"`
	}
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &translated); err != nil {
		return nil, fmt.Errorf("decode translations: %w", err)
	}
	if len(translated.Translations) != len(text) {
		return nil, fmt.Errorf("expected %d translations, got %d", len(text), len(translated.Translations))
	}

	return translated.Translations, nil
}
