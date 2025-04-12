package main

import "fmt"

func main() {
	arr1 := [3]int{1, 2, 3}
	slc1 := []int{1, 2, 3}

	fmt.Println(arr1)
	fmt.Println(slc1)
	// Por padrão, o Go tem apenas duas estruturas de "coleção" por índice, o array e o slice:

	// Um slice em go funciona como um ArrayList em Java

	// - Capacidade dinâmica, uma abstração do array
	// - Sua capacidade pode ser inicializada com o valor que quiser, mas é bom que já seja iniciado com o valor pretendido
	// - Sempre que atinge sua capacidade (valor que pode ser resgatado com a função `cap()`), ela é dobrada (criando um novo array na memória em uma diferente posição, e deletando o antigo)

	// Lembretes gerais:
	//
	// - Sempre que se cria um novo subset de um array, o que se faz é registrar um novo ponteiro para os mesmos endereços da memória. Caso queira uma cópia, use `copy(dest, src)`
	// - Obviamente um array nunca diminui de tamanho, apenas aumenta (no caso da abstração do slice, que não aumenta literalmente, mas você entendeu)
	slc2 := []int{1, 2, 3}

	fmt.Println(slc2)
	fmt.Println("Capacity:", cap(slc2))
	fmt.Println("Length:", len(slc2))

	slc2 = append(slc2, 4)

	fmt.Println("Capacity after append:", cap(slc2))
	fmt.Println("Length after append:", len(slc2))

	// Curiosidade: Os métodos `len()` e `cap()` também podem ser utilizados em ponteiros de arrays, mas não em ponteiros de slices.

	fmt.Println("Array pointer length:", len(&arr1))
	fmt.Println("Array pointer capacity:", cap(&arr1))

	// O único método de iteração de um array é o `for range`, que retorna o índice e o valor do elemento. O índice pode ser omitido, mas o valor não pode.
	// - O `for range` não retorna o endereço de memória, mas sim o valor do elemento. Para pegar o endereço, use `&array[i]`
	// - Ao utilizar `&v`, o valor retornado é o endereço de memória do valor presente apenas no momento da iteração, que é diferente do endereço de memória do valor que está na posição `i` do array.

	fmt.Println("\nForRange demonstration:")
	for i, v := range slc2 {
		fmt.Printf("Index: %d, Value: %v\n", i, v)
		fmt.Printf("Value address: %p\n", &v)
		fmt.Printf("Slice element address: %p\n\n", &slc2[i])
	}
}
