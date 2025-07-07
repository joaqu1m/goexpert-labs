package main

func main() {

	ch := make(chan int)

	go publish(ch)
	consume(ch)

}

func publish(ch chan<- int) { // canal de escrita, parâmetro definido como read-only
	for i := range 10 {
		ch <- i
	}
	close(ch)
}

func consume(ch <-chan int) { // canal de leitura, parâmetro definido como write-only
	for x := range ch {
		println("Received:", x)
	}
}
