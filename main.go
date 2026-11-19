package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/translate", returnTranslateHandler)
	mux.HandleFunc("/api/inventory", returnInventoryValueHandler)

	err := http.ListenAndServe(":3000", mux)
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
