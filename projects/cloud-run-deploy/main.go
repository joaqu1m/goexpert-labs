package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	viacepBase  = "https://viacep.com.br/ws"
	weatherBase = "http://api.weatherapi.com/v1/current.json"
)

var weatherKey string

func main() {
	_ = godotenv.Load()

	weatherKey = os.Getenv("WEATHERAPI_KEY")
	if weatherKey == "" {
		panic("WEATHERAPI_KEY não definido")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/weather", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !regexp.MustCompile(`^\d{8}$`).MatchString(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s/json/", viacepBase, cep))
	if err != nil {
		http.Error(w, "error viacep", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	var via struct {
		Localidade string `json:"localidade"`
		UF         string `json:"uf"`
		Erro       bool   `json:"erro"`
	}
	json.NewDecoder(resp.Body).Decode(&via)
	if via.Erro || via.Localidade == "" {
		// não seria could?
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	q := url.QueryEscape(fmt.Sprintf("%s,%s,BR", via.Localidade, via.UF))
	weatherURL := fmt.Sprintf("%s?key=%s&q=%s", weatherBase, weatherKey, q)
	resp, err = http.Get(weatherURL)
	if err != nil {
		http.Error(w, "error weatherapi", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	var wea struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	json.NewDecoder(resp.Body).Decode(&wea)

	tempC := wea.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273
	out := map[string]string{
		"temp_C": strconv.FormatFloat(tempC, 'f', 1, 64),
		"temp_F": strconv.FormatFloat(tempF, 'f', 1, 64),
		"temp_K": strconv.FormatFloat(tempK, 'f', 1, 64),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
