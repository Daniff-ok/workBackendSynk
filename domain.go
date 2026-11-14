package main

var Apikey string

// INVENTORY VALUE
type InventoryRequest struct {
	Token      string          `json:"token"`
	Hwid       string          `json:"hwid"`
	HwidStable string          `json:"hwidStable"`
	Game       int             `json:"game"`
	Players    []requestPlayer `json:"items"`
}

type requestPlayer struct {
	ID        string `json:"steamId"`
	EnemyTeam bool   `json:"team2"`
}

type ReturnInventoryValue struct {
	Ready   bool
	Players []returnplayer
}

type returnplayer struct {
	ID        string
	Name      string
	value     float64
	Locked    float64
	Blocked   float64
	busy      bool
	avatar    string
	EnemyTeam bool
	Party     int
}

// TRANSLATE
type TRrequest struct {
	Token      string  `json:"token"`
	Hwid       string  `json:"hwid"`
	HwidStable string  `json:"hwidStable"`
	TargetLang string  `json:"targetLang"`
	Game       int     `json:"game"`
	Items      []items `json:"items"`
}

type items struct {
	Text           string `json:"text"`
	SourceLanguage string `json:"sourceLang"`
}

type translateText struct {
	Text   string `json:"text"`
	Source string `json:"scr"`
}

type TranslateReturn struct {
	Results []translateText `json:"results"`
}
