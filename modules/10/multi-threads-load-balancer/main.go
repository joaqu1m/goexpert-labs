package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan int)

	workers := 2_000

	for i := range workers {
		go consume(i, ch)
	}
	publish(ch)

}

func publish(ch chan<- int) { // canal de escrita, parâmetro definido como read-only
	for i := range 10_000 {
		ch <- i
	}
}

// aqui o consume atua como um worker, recebendo os dados do canal
func consume(workerID int, ch <-chan int) {
	for x := range ch {
		fmt.Printf("Worker %d is processing: %d\n", workerID, x)
		// simula processamento
		time.Sleep(time.Second)
	}
}
