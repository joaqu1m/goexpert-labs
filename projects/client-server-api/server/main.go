package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	db, err := sql.Open("sqlite3", "cotacoes.db")
	if err != nil {
		log.Fatalf("falha ao abrir banco: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		log.Fatalf("falha ao criar tabela: %v", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		apiCtx, cancelAPI := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancelAPI()

		req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, "https://economia.awesomeapi.com.br/json/usd/", nil)
		if err != nil {
			log.Printf("erro ao criar requisição: %v", err)
			http.Error(w, "erro interno", http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if apiCtx.Err() == context.DeadlineExceeded {
				log.Printf("timeout na chamada da API: %v", err)
			} else {
				log.Printf("erro ao chamar API: %v", err)
			}
			http.Error(w, "erro ao obter cotação", http.StatusGatewayTimeout)
			return
		}
		defer resp.Body.Close()

		var data []Cotacao
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("erro ao decodificar resposta da API: %v", err)
			http.Error(w, "erro interno", http.StatusInternalServerError)
			return
		}
		if len(data) == 0 {
			log.Printf("resposta vazia da API")
			http.Error(w, "nenhuma cotação disponível", http.StatusNoContent)
			return
		}

		cot := data[0]

		dbCtx, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelDB()

		_, err = db.ExecContext(dbCtx, "INSERT INTO quotes (bid) VALUES (?)", cot.Bid)
		if err != nil {
			if dbCtx.Err() == context.DeadlineExceeded {
				log.Printf("timeout ao persistir no banco: %v", err)
			} else {
				log.Printf("erro ao inserir no banco: %v", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cot.Bid); err != nil {
			log.Printf("erro ao enviar resposta: %v", err)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("erro no servidor HTTP: %v", err)
	}
}
