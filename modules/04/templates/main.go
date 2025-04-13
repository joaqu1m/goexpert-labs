package main

import (
	// use "text/template" quando não for usar HTML
	// o "html/template" tem algumas funções a mais para evitar XSS
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Pessoa struct {
	Nome  string
	Idade int
}

type Pedido struct {
	ID         int
	Produto    string
	Quantidade int
	Valor      float64
	Data       string
}

type Pedidos struct {
	Pessoa  Pessoa
	Pedidos []Pedido
}

func main() {
	pessoa := Pessoa{Nome: "John Doe", Idade: 30}

	testTemplate := template.Must(
		template.New("example").Parse("Hello, {{.Nome}}! You are {{.Idade}} years old.\n"),
	)

	err := testTemplate.ExecuteTemplate(os.Stdout, "example", pessoa)
	if err != nil {
		panic(err)
	}

	pedidos := Pedidos{
		Pessoa: pessoa,
		Pedidos: []Pedido{
			{ID: 1, Produto: "Fralda", Quantidade: 2, Valor: 10.50, Data: "2023-10-01"},
			{ID: 2, Produto: "Picanha", Quantidade: 1, Valor: 20.00, Data: "2023-10-02"},
			{ID: 3, Produto: "Cerveja", Quantidade: 5, Valor: 5.00, Data: "2023-10-03"},
		},
	}

	// Lembrando que o método ParseFiles aceita vários arquivos, o que podem servir como uma espécie de "diferentes componentes"; Por exemplo aquele "header.html" ali embaixo
	// O nome do template sempre deverá ser o nome do arquivo principal, e por dentro do arquivo principal definimos onde estarão os outros com {{ template "X" }}
	realTemplate := template.New("pedidos.html")

	// bora colocar um pouco de lógica nisso ae!
	realTemplate.Funcs(template.FuncMap{
		"Upper": strings.ToUpper,
	})

	realTemplate, err = realTemplate.ParseFiles("pedidos.html", "header.html")
	if err != nil {
		panic(err)
	}

	// bora hospedar isso ae num server!
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = realTemplate.Execute(w, pedidos)
		if err != nil {
			panic(err)
		}
	})
	http.ListenAndServe(":8080", nil)

}
