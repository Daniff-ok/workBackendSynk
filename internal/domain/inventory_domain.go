package domain

// INVENTORY VALUE
type InventoryRequest struct {
	Token      string          `json:"token"`
	Hwid       string          `json:"hwid"`
	HwidStable string          `json:"hwidStable"`
	Game       int             `json:"game"`
	Players    []RequestPlayer `json:"roster"`
}

type RequestPlayer struct {
	ID        string `json:"steamId"`
	EnemyTeam bool   `json:"team2"`
}

type ReturnInventoryValue struct {
	Ready     bool                `json:"ready"`
	Players   []ReturnPlayer      `json:"players"`
	GameValue map[int][]GameValue `json:"byGame"`
}

type ReturnPlayer struct {
	ID        string   `json:"steamId"`
	Name      string   `json:"name"`
	Value     *float64 `json:"value"`
	Avatar    string   `json:"avatar"`
	EnemyTeam bool     `json:"team2"`
	Party     int      `json:"party"`
}

type GameValue struct {
	ID    string  `json:"steamId"`
	Value float64 `json:"value"`
}

// API
type ButchWokRequest struct {
	SteamIds []string `json:"steamids"`
	Games    []string `json:"games"`
}

type WokResponse struct {
	Results []struct {
		ID         string  `json:"steamid"`
		Game       string  `json:"game"`
		TotalValue float64 `json:"total"`
	} `json:"results"`
}

type SteamResponce struct {
	Resp struct {
		Players []PlayerInfo `json:"players"`
	} `json:"response"`
}

type PlayerInfo struct {
	ID     string `json:"steamid"`
	Name   string `json:"personaname"`
	Avatar string `json:"avatar"`
}
