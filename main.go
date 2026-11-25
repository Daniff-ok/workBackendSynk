package main

import (
	"log/slog"
	"net/http"
	"os"

	"rndmBackend/internal/handler"
	"rndmBackend/internal/service"
	"rndmBackend/internal/gateway"
	"rndmBackend/internal/repo"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	db := GetDbConnection()

	loginHandler := handler.NewLoginHandler(service.NewLoginService(gateway.NewLogingateway(),repo.NewLoginRepo(db)))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/login", loginHandler.Login)
	mux.HandleFunc("GET /api/refresh", loginHandler.Refresh)

	mux.HandleFunc("/api/translate", returnTranslateHandler)
	mux.HandleFunc("/api/inventory", returnInventoryValueHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("Server starting", "error: ", err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Ошибка загрузки .env файла")
	}
}
