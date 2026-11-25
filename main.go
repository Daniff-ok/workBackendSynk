package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
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

	mux.HandleFunc("/api/translate", returnTranslateHandler)
	mux.HandleFunc("/api/inventory", returnInventoryValueHandler)
	mux.HandleFunc("/", dumpHandler)


	err := http.ListenAndServe(":3001", mux)
	if err != nil {
		slog.Error("Server starting", "error: ", err)
	}
}

func dumpHandler(w http.ResponseWriter, r *http.Request) {

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("%s\n", dump)

	w.Write([]byte(""))
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Ошибка загрузки .env файла")
	}
}
