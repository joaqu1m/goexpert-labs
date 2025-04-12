package main

import "fmt"

func main() {
	hashmap := make(map[string]int)
	hashmap2 := map[string]int{}

	fmt.Println(hashmap)
	fmt.Println(hashmap2)

	// Diferenças entre os dois métodos de criação de hashmap:
	// - `make()`: Pode ter um tamanho inicial especificado, para poupar esforço posterior de redimensionamento.
	// - Com um composite literal: O tamanho inicial é sempre 0, mas pode facilitar a inserção de elementos no momento da criação.
	//
	// Apesar dessas pequenas diferenças, ambos fazem a exata mesma coisa.

	hashmap["key"] = 1
	hashmap["key2"] = 2
	hashmap["key3"] = 3

	for key, value := range hashmap {
		fmt.Println(key, value)
	}

	delete(hashmap, "key2")
	fmt.Println(hashmap)

	// O `delete()` muta diretamente o hashmap original.
}
