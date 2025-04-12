package main

import (
	"fmt"
	"strings"
)

func main() {
	// Basic function example
	a := []int{1, 2, 3}

	fmt.Println(sum(a))

	// Naturalmente, caso não seja feita uma cópia, funções utilizam a exata mesma instância do objeto,
	// o que pode levar a resultados inesperados.
	a2 := []int{1, 2, 3}
	fmt.Println(&a2[0])
	sumWithPrint(a2)

	// Demonstrating that slices are passed by reference
	b := []int{1, 2, 3}
	fmt.Println(doubleAll(b))
	fmt.Println(b)

	// Funções variádicas em go - mesma coisa que em qualquer outra lang, ou ao spread operator do js
	print(mergeln("a", "b", "c"))
	print(mergeln("a", "b", "c"))

	// Closure:
	// - Função aninhada que tem acesso ao escopo da função pai
	// - Pode ser uma opção mais rápida para criar uma função que não precisa ser reutilizada
	//
	// Se assemelha à funções lambda ou funções anônimas de outras linguagens
	duplicar := multiplicadorPor(2)
	triplicar := multiplicadorPor(3)

	fmt.Println("5 duplicado:", duplicar(5))
	fmt.Println("5 triplicado:", triplicar(5))
	fmt.Println("10 duplicado:", duplicar(10))
}

func sum(a []int) int {
	total := 0
	for _, v := range a {
		total += v
	}
	return total
}

func sumWithPrint(a []int) int {
	fmt.Println(&a[0])
	total := 0
	for _, v := range a {
		total += v
	}
	return total
}

func doubleAll(arr []int) []int {
	for i := range arr {
		arr[i] *= 2
	}
	return arr
}

func mergeln(args ...string) string {
	return strings.Join(args, " ") + "\n"
}

func multiplicadorPor(fator int) func(int) int {
	return func(valor int) int {
		return valor * fator
	}
}
