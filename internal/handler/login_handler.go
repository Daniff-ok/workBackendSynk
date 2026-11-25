package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"rndmBackend/internal/service"
)

type LoginHandler struct {
	service *service.LoginService
}

func NewLoginHandler(service *service.LoginService) *LoginHandler {
	return &LoginHandler{service}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	code := ""
	err := json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		slog.Error("Code json decode error", "error", err)
		http.Error(w, "Invalid code", http.StatusBadRequest)
	}

	r.Context()
	responceStruct, err := h.service.Login(code, r.Context())
	if err != nil {
		slog.Error("Login error", "err", err)
		http.Error(w, "Login error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(responceStruct)
	if err != nil {
		slog.Error("Json Encode error", "err", err)
		http.Error(w, "Json Encode error", http.StatusInternalServerError)
		return
	}
}

func (h *LoginHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	token := ""
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		slog.Error("Code json decode error", "error", err)
		http.Error(w, "Invalid token", http.StatusBadRequest)
	}

	r.Context()
	responceStruct, err := h.service.Refresh(token, r.Context())
	if err != nil {
		slog.Error("Refresh error", "err", err)
		http.Error(w, "Login error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(responceStruct)
	if err != nil {
		slog.Error("Json Encode error", "err", err)
		http.Error(w, "Json Encode error", http.StatusInternalServerError)
		return
	}
}
