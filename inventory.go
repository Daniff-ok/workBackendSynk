package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var wokApiUrl = "https://woksteamapi.com/v1/inventories"
var steamApiUrl = "https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002"

var inventoryHTTP = &http.Client{Timeout: 3 * time.Minute}

var inventoryGames = map[string]int{"dota2": 570, "cs2": 730, "rust": 252490}

func returnInventoryValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	var inventory InventoryRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32_768)).Decode(&inventory); err != nil ||
		len(inventory.Players) == 0 || len(inventory.Players) > 10 {
		http.Error(w, "Expected a roster of 1 to 10 players", http.StatusBadRequest)
		return
	}

	request := butchWokRequest{Games: []string{"dota2", "cs2", "rust"}}
	seen := make(map[string]bool, len(inventory.Players))
	for _, player := range inventory.Players {
		const individualBase = 76561197960265728
		id, err := strconv.ParseUint(player.ID, 10, 64)
		if err != nil || id <= individualBase || id > individualBase+0xffffffff || seen[player.ID] ||
			strconv.FormatUint(id, 10) != player.ID {
			http.Error(w, "Invalid or duplicate Steam ID", http.StatusBadRequest)
			return
		}
		seen[player.ID] = true
		request.SteamIds = append(request.SteamIds, player.ID)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()
	values, err := requestInventory(ctx, request, inventoryHTTP)
	if err != nil {
		slog.Warn("Inventory provider request failed", "error", err)
		http.Error(w, "Inventory provider unavailable", http.StatusBadGateway)
		return
	}

	profiles, err := requestAddInfo(ctx, request.SteamIds, inventoryHTTP)
	if err != nil {
		slog.Warn("Steam profile enrichment failed")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(mapInventory(values, profiles, inventory)); err != nil {
		slog.Warn("Inventory response failed", "error", err)
	}
}

func requestInventory(ctx context.Context, requestBody butchWokRequest, cl *http.Client) (wokResponse, error) {
	data, err := json.Marshal(requestBody)
	if err != nil {
		return wokResponse{}, fmt.Errorf("Json marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, wokApiUrl, bytes.NewReader(data))
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

	if resp.Status != "200 OK" {
		return wokResponse{}, fmt.Errorf("Return error from inventory api, status %v, request body: %v", resp.Status, requestBody)
	}

	var dataResp wokResponse
	err = json.NewDecoder(resp.Body).Decode(&dataResp)
	if err != nil {
		return wokResponse{}, fmt.Errorf("Json decode error: %w", err)
	}

	return dataResp, nil
}

func requestAddInfo(ctx context.Context, steamids []string, cl *http.Client) (steamResponce, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, steamApiUrl, nil)
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

	if resp.Status != "200 OK" {
		return steamResponce{}, fmt.Errorf("Steam profiles returned %s", resp.Status)
	}

	var dataResp steamResponce
	err = json.NewDecoder(resp.Body).Decode(&dataResp)

	return dataResp, err
}

func mapInventory(values wokResponse, profiles steamResponce, roster InventoryRequest) ReturnInventoryValue {
	response := ReturnInventoryValue{
		Ready:     true,
		Players:   make([]returnPlayer, len(roster.Players)),
		GameValue: make(map[int][]gameValue),
	}
	profileByID := make(map[string]PlayerInfo, len(profiles.Resp.Players))
	for _, profile := range profiles.Resp.Players {
		profileByID[profile.ID] = profile
	}
	playerIndex := make(map[string]int, len(roster.Players))
	for i, player := range roster.Players {
		profile := profileByID[player.ID]
		response.Players[i] = returnPlayer{
			ID: player.ID, Name: profile.Name, Avatar: profile.Avatar, EnemyTeam: player.EnemyTeam,
		}
		playerIndex[player.ID] = i
	}
	valuations := make(map[int]map[string]float64)
	for _, result := range values.Results {
		appID, gameKnown := inventoryGames[result.Game]
		_, playerKnown := playerIndex[result.ID]
		if !gameKnown || !playerKnown {
			continue
		}
		if valuations[appID] == nil {
			valuations[appID] = make(map[string]float64)
		}
		valuations[appID][result.ID] = result.TotalValue
	}
	for appID, gameValues := range valuations {
		for id, value := range gameValues {
			response.GameValue[appID] = append(response.GameValue[appID], gameValue{ID: id, Value: value})
			player := &response.Players[playerIndex[id]]
			if player.Value == nil {
				player.Value = new(float64)
			}
			*player.Value += value
		}
	}
	return response
}
