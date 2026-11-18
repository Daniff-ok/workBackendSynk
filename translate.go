package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	var requestData TRrequest

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Json decode error", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	text := []string{}
	for _, val := range requestData.Items {
		text = append(text, val.Text)
	}

	if len(text) >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("гембит не роби"))
		return
	}

	trReturn := TranslateReturn{}
	translate, err := Translate(text, "ru")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("deepseek error"))
		return
	}

	for _, val := range translate {
		trReturn.Results = append(trReturn.Results, TranslateText{val, "zh"})
	}

	fmt.Println(trReturn.Results)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(trReturn)
	if err != nil {
		fmt.Println(err)
	}
}

func Translate(text []string, language string) ([]string, error) {
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

	req, err := http.NewRequest(
		"POST",
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
