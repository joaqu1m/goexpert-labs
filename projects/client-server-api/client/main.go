package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	apiCtx, cancelAPI := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancelAPI()

	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Printf("erro ao criar requisição: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if apiCtx.Err() == context.DeadlineExceeded {
			log.Printf("timeout na chamada da API: %v", err)
		} else {
			log.Printf("erro ao chamar API: %v", err)
		}
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status da resposta:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		log.Printf("erro na resposta da API: %s", resp.Status)
		return
	}

	var cotacao string
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		log.Printf("erro ao decodificar resposta: %v", err)
		return
	}
	fmt.Println("Cotação obtida:", cotacao)
}
