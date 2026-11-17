package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
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
	mux.HandleFunc("/", dumpHandler)


	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		fmt.Println(err)
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

func returnTranslateHandler(w http.ResponseWriter, r *http.Request) {

	var requestData TRrequest
	fmt.Println("Get translate")

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Json error"))
		return
	}

	defer r.Body.Close()

	text := []string{}
	for _, val := range requestData.Items {
		text = append(text, val.Text)
	}

	if len(text) >= 10 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("гембит не роби"))
		return
	}

	trReturn := TranslateReturn{}
	translate, err := Translate(text, "ru")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("deepseek error"))
		return
	}

	for _, val := range translate {
		trReturn.Results = append(trReturn.Results, TranslateText{val, "zh"})
	}

	fmt.Println(trReturn.Results)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(trReturn)
	if err != nil {
		fmt.Println(err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки .env файла")
	}
}
