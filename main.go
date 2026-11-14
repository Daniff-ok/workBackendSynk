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

var apikey string

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

type translateText struct {
	Text   string `json:"text"`
	Source string `json:"scr"`
}

type TranslateReturn struct {
	Results []translateText `json:"results"`
}

func main() {
	loadEnv()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/translate", returnTranslateHandler)
	mux.HandleFunc("/", dumpHandler)

	err := http.ListenAndServe(":6767", mux)
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

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Json error"))
		return
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

		trReturn.Results = append(trReturn.Results, translateText{translated, "zh"})

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("translate error"))
		}
	}

	fmt.Println(trReturn.Results)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

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

	apikey = os.Getenv("STEAM_INVENTORY_API")
}
