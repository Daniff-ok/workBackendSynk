package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/bregydoc/gtranslate"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/translate", returnTranslateHandler)
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
	fmt.Println("Get request")
	fmt.Println(requestData)

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Json error"))
	}

	var trReturn TranslateReturn

	for _, val := range requestData.Items {
		translated, err := gtranslate.TranslateWithParams(
			val.Text,
			gtranslate.TranslationParams{
				From: "zh",
				To:   "ru",
			},
		)

		trReturn.Results = append(trReturn.Results, TranslateText{translated, "zh"})

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("translate error"))

		}
	}

	fmt.Println(trReturn.Results)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	fmt.Println(trReturn)
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

	Apikey = os.Getenv("STEAM_INVENTORY_API")
}
