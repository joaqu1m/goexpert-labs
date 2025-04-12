package main

func Soma[T int | float32 | float64](a, b T) T {
	return a + b
}

type number = interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// Soma implementando "constraints", um tipo que pode representar vários outros diferentes tipos
func Soma2[T number](a, b T) T {
	return a + b
}

// Constraints geralmente restringem o tipo a ponto de não conseguir implementar nem instâncias de si mesmo. exemplo:

type Inteiro int

// Para fazer com que "Inteiro" seja um tipo válido, é necessário permitir todos os tipos de "int" usando a notação "~":
type number2 = interface {
	~int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// Existem constraints pré-definidas em Go, como "comparable", "constraints.Ordered" e "constraints.Complex", que podem ser usadas para restringir os tipos a números inteiros, números de ponto flutuante ou números complexos. Essas constraints são úteis quando você deseja garantir que um tipo genérico seja um número.

func Comparador[T comparable](a, b T) bool {
	return a == b
}

func main() {
	// Exemplo de uso da função Soma com diferentes tipos
	intResult := Soma(10, 20)
	floatResult := Soma(10.5, 20.5)

	println("Resultado da soma de inteiros:", intResult)                                    // Resultado: 30
	println("Resultado da soma de floats:", floatResult)                                    // Resultado: 31
	println("Resultado da soma de floats com float32:", Soma(float32(10.5), float32(20.5))) // Resultado: 31

	println(Comparador(10, 20))
}
