package main



// INVENTORY VALUE
type InventoryRequest struct {
	Token      string          `json:"token"`
	Hwid       string          `json:"hwid"`
	HwidStable string          `json:"hwidStable"`
	Game       int             `json:"game"`
	Players    []requestPlayer `json:"roster"`
}

type requestPlayer struct {
	ID        string `json:"steamId"`
	EnemyTeam bool   `json:"team2"`
}

type ReturnInventoryValue struct {
	Ready     bool                 `json:"ready"`
	Players   [9]returnPlayer      `json:"players"`
	GameValue map[string][]gameValue `json:"byGame"`
}

type returnPlayer struct {
	ID        string  `json:"steamid"`
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Avatar    string  `json:"avatar"`
	EnemyTeam bool    `json:"team2"`
	Party     int     `json:"party"`
}

type gameValue struct {
	ID    string  `json:"steamId"`
	Value float64 `json:"value"`
}



// wok, steam api
type butchWokRequest struct {
	SteamIds []string `json:"steamids"`
	Games    []string `json:"games"`
}

type wokResponse struct {
	Results []struct {
		ID         string `json:"steamid"`
		Game       string `json:"game"`
		TotalValue float64 `json:"total"`
	} `json:"results"`
}

type steamResponce struct {
	Resp struct {
		Players []PlayerInfo `json:"players"`
	} `json:"response"`
}

type PlayerInfo struct {
	ID     string `json:"steamid"`
	Name   string `json:"personaname"`
	Avatar string `json:"avatar"`
}



// TRANSLATE
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
	Source string `json:"scr"`
}

type TranslateReturn struct {
	Results []TranslateText `json:"results"`
}
