package main

import (
	"fmt"
	"net/http"
)

// mux: Multiplexer (ou multiplexador), usado em todas as langs para implementar servidores http, porém explícito quando se trata do go
func main() {
	// instanciar um mux ao em vez de utilizar o padrão possibilita o uso de diferentes portas dentro do mesmo servidor (seja lá qual for o motivo de fazer algo assim)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// também é possível criar uma struct para servir diretamente como handler (ou um endpoint), e com isso se torna possível passar parâmetros para o handler de um modo mais fácil
	mux.Handle("/hello", &MyHandler{
		Param1: "param1",
		Param2: 42,
	})

	http.ListenAndServe(":8080", mux)
}

type MyHandler struct {
	Param1 string
	Param2 int
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hello, %s! You are %d years old.", h.Param1, h.Param2)))
}
