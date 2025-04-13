package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// Também é possível adicionar valores ao contexto
	// O valor é passado como um par chave-valor
	// O sentido disso é poder passar informações entre funções que possuem o mesmo contexto de maneira mais fácil
	ctx = context.WithValue(ctx, "key", "value")

	doSomething(ctx)
}

func doSomething(ctx context.Context) {
	fmt.Println("Starting work")
	fmt.Println("Value from context:", ctx.Value("key"))
	// O statement "select" itera infinitamente até que um dos casos seja atendido
	// O caso "ctx.Done()" é chamado quando o contexto é cancelado ou expira... e nesse caso, ele sempre expira
	select {
	case <-time.After(time.Second * 10):
		fmt.Println("Work done")
	case <-ctx.Done():
		fmt.Println("Context cancelled or timed out")
	}
	fmt.Println("Exiting function")
}
