package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

// O context também pode ser adquirido do request em um handler
// O contexto é cancelado quando o request é cancelado ou expira
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Neste exemplo eu espero por 5 segundos para simular um trabalho
	// O contexto é cancelado quando o request é cancelado ou expira antes desses 5 segundos
	log.Println("Starting request")
	select {
	case <-ctx.Done():
		log.Println("Request cancelled")
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		return
	case <-time.After(5 * time.Second):
		c := 0
		for i := range 1_000_000 {
			log.Println("Processing request", i)
			c++
		}
		log.Println("Total processed requests:", c)
		w.Write([]byte("Request completed"))
	}
}
