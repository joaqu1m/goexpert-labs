package main

func main() {
	var x interface{} = 10

	if v, ok := x.(int); ok {
		println("x é um inteiro:", v)
	} else {
		println("x não é um inteiro")
	}

	// type any = interface{}
	// apenas um alias moderno para interface{}
	var y any = "Hello"
	if str, ok := y.(string); ok {
		println("y é uma string:", str)
	} else {
		println("y não é uma string")
	}

	var z any = 3.14
	v := z.(float64)
	println("z é um float64:", v)
	// v := z.(int) // Isso causaria um panic, pois z não é um int
	// v, ok := z.(int) // Isso não causaria um panic, mas v seria zero

}
