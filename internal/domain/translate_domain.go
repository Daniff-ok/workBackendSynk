package domain

type TRrequest struct {
	Token      string  `json:"token"`
	Hwid       string  `json:"hwid"`
	HwidStable string  `json:"hwidStable"`
	TargetLang string  `json:"targetLang"`
	Game       int     `json:"game"`
	Items      []Items `json:"items"`
}

type Items struct {
	Text           string `json:"text"`
	SourceLanguage string `json:"sourceLang"`
}

type TranslateText struct {
	Text   string `json:"text"`
	Source string `json:"src"`
}

type TranslateReturn struct {
	Results []TranslateText `json:"results"`
}

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