package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func fetchAndChannelize(url string, ch chan<- map[string]any) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return
	}

	ch <- body
}

func main() {

	cep := "01001000"

	ch1 := make(chan map[string]any)
	ch2 := make(chan map[string]any)

	go fetchAndChannelize(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), ch1)
	go fetchAndChannelize(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep), ch2)

	select {
	case result := <-ch1:
		fmt.Println("Response from BrasilAPI:", result)
	case result := <-ch2:
		fmt.Println("Response from ViaCEP:", result)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: No response received within the specified duration")
	}
}
