package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// aqui começamos a utilizar "tags", anotações que podem ser utilizadas por qualquer um que consultar a struct.
// No caso, estamos utilizando a tag "json" pois as funções da biblioteca json
// utilizam essa tag para identificar como os dados devem ser serializados e deserializados

// informações que não são necessárias no json podem ser omitidas da serialização com "-"
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Mood string `json:"-"`
}

func main() {

	person := Person{Name: "John Doe", Age: 30, Mood: "happy"}
	fmt.Println(person)

	jsonData, err := json.Marshal(person)
	if err != nil {
		fmt.Println("Error serializing JSON:", err)
		return
	}
	fmt.Println(string(jsonData))

	// para serializar e imediatamente utilizar o recurso, ao em vez de salvar na memória:
	err = json.NewEncoder(os.Stdout).Encode(person)
	if err != nil {
		fmt.Println("Error serializing JSON:", err)
		return
	}

	jsonData2 := []byte(`{"name":"Jane Doe","age":25, "Mood":"sad"}`)
	// para deserializar
	var p Person
	err = json.Unmarshal(jsonData2, &p)
	if err != nil {
		fmt.Println("Error deserializing JSON:", err)
		return
	}
	fmt.Println(p)
}
