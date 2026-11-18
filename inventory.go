package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"strings"
)

var wokApiUrl = "https://woksteamapi.com/v1/inventories"
var steamApiUrl = "https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002"

func returnInventoryValueHandler(w http.ResponseWriter, r *http.Request) {

	var inventoryRequestData InventoryRequest
	slog.Debug("Get request")

	err := json.NewDecoder(r.Body).Decode(&inventoryRequestData)
	if err != nil {
		slog.Error("Json decode error", "err", err, "status", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Json decode request error"))
		return
	}

	defer r.Body.Close()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	request := butchWokRequest{}
	for _, player := range inventoryRequestData.Players{
		request.SteamIds = append(request.SteamIds, player.ID)
	}
	request.Games = []string{"dota2", "cs2"} // add rust

	wokResp, err := requestInventory(request, client)
	if err != nil {
		slog.Error("Request Inventory error", "err", err)
		http.Error(w, "Request Inventory from api error", http.StatusBadGateway)
		return
	}

	steamRest, err := requestAddInfo(request.SteamIds, client)
	if err != nil {
		slog.Error("Request Inventory error", "err", err)
		http.Error(w, "Request steam addititon information from api error", http.StatusBadGateway)
		return
	}
	
	mappedData := maptoAppResponce(wokResp, steamRest, inventoryRequestData)


	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(mappedData)
	if err != nil {
		slog.Error("Json Encode error", "err", err)
		http.Error(w, "Json Encode error", http.StatusInternalServerError)
		return
	}

	slog.Debug("Succefull sent responce")
}

func requestInventory(requestBody butchWokRequest, cl *http.Client) (wokResponse, error) {
	data, err := json.Marshal(requestBody)
	if err != nil {
		return wokResponse{}, fmt.Errorf("Json marshal error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, wokApiUrl, bytes.NewReader(data))
	if err != nil {
		return wokResponse{}, fmt.Errorf("Create new request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("STEAM_INVENTORY_API"))

	slog.Debug("Do wok inventory butch request")
	resp, err := cl.Do(req)
	if err != nil {
		slog.Error("Connection error to inventory checker", "api", wokApiUrl, "status", http.StatusBadGateway)
		return wokResponse{}, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK"{
		return wokResponse{}, fmt.Errorf("Return error from inventory api, status %v, request body: %v", resp.Status, requestBody)
	}

	var dataResp wokResponse
	err = json.NewDecoder(resp.Body).Decode(&dataResp)
	if err != nil{
		return wokResponse{}, fmt.Errorf("Json decode error: %w", err)
	}

	return dataResp, nil
}

func requestAddInfo(steamids []string, cl *http.Client) (steamResponce, error){
	req, err := http.NewRequest(http.MethodGet, steamApiUrl, nil)
	if err != nil {
		return steamResponce{}, fmt.Errorf("Create new request error: %w", err)
	}

	q := req.URL.Query()
	q.Add("key", os.Getenv("STEAM_WEB_API"))
	q.Add("steamids", strings.Join(steamids, ","))
	req.URL.RawQuery = q.Encode()

	resp, err := cl.Do(req)
	if err != nil {
		return steamResponce{}, fmt.Errorf("Connection error to inventory checker %w, status %d", err, http.StatusBadGateway)
	}

	slog.Debug("Do steam official api add info request")
	defer resp.Body.Close()

	if resp.Status != "200 OK"{
		return steamResponce{}, fmt.Errorf("Return error from steam official api, status %v, ulr: %v", resp.Status, req.URL)
	}

	var dataResp steamResponce
	err = json.NewDecoder(resp.Body).Decode(&dataResp)

	return dataResp, nil
}

func maptoAppResponce(wokResp wokResponse, steamResp steamResponce, reqData InventoryRequest) ReturnInventoryValue {
	var returnInventoryValue ReturnInventoryValue
	returnInventoryValue.Ready = true

	for i, player := range reqData.Players {
		for _, infoPlayer:= range steamResp.Resp.Players {
			if player.ID == infoPlayer.ID {
				returnInventoryValue.Players[i] = returnPlayer{
					ID: player.ID,
					Name: infoPlayer.Name,
					Avatar: infoPlayer.Avatar,
					EnemyTeam: player.EnemyTeam,
					Party: 0,
				}
			}
		}
	}

	returnInventoryValue.GameValue = make(map[string]gameValue)
	for i, player := range returnInventoryValue.Players {
		for _, pairPlayerGame := range wokResp.Results {
			if player.ID == pairPlayerGame.ID {
				returnInventoryValue.Players[i].Value += pairPlayerGame.TotalValue

				returnInventoryValue.GameValue[pairPlayerGame.Game] = gameValue{pairPlayerGame.ID, pairPlayerGame.TotalValue}
			}
		}
	}

	return returnInventoryValue
}

func wokApiPlaceholder(_ butchWokRequest, _ *http.Client) (wokResponse, error) {
	return wokResponse{
		Results: []struct{
			ID         string `json:"steamid"`
			Game       string `json:"game"`
			TotalValue float64 `json:"total"`
		}{
			{
				ID: "76561198147927273",
				Game: "dota2",
				TotalValue: 595.31,
			},
			{
				ID: "76561198147927273",
				Game: "cs2",
				TotalValue: 105,
			},
						{
				ID: "76561198063639479",
				Game: "dota2",
				TotalValue: 780.99,
			},
						{
				ID: "76561198063639479",
				Game: "cs2",
				TotalValue: 5.5,
			},
		},
	}, nil
}