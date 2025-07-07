package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var visitors = 0

func main() {

	m := sync.Mutex{}

	// problema clássico de concorrência:
	// cada requisição HTTP abre uma nova thread
	// e todas elas compartilham a variável visitors
	// se duas requisições forem feitas ao mesmo tempo,
	// elas podem ler o mesmo valor de visitors
	// e incrementar o mesmo valor, resultando em um número errado de visitantes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m.Lock() // bloqueia o mutex para garantir que apenas uma thread/goroutine acesse a variável visitors por vez
		visitors++
		m.Unlock() // libera o mutex para que outras threads/goroutines possam acessar a variável visitors
		println("Visitor number:", visitors)

		// uma opção mais simples seria usar o pacote atomic
		// atomic.AddInt32(&visitors, 1)
		// mas vamos usar o mutex para entender como funciona a concorrência
		// na prática, o mutex é mais flexível e pode ser usado para proteger outras variáveis ou estruturas de dados
		// e o pacote atomic é mais simples e rápido para operações atômicas em variáveis inteiras
		// mas eles fazem a mesma coisa

		time.Sleep(2 * time.Second)

		w.Write([]byte(fmt.Sprintf("Hello, visitor number %d!", visitors)))
	})
	http.ListenAndServe(":8000", nil)
}
