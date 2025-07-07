package main

func main() {

	ch := make(chan int)

	go publish(ch)
	consume(ch)

}

func publish(ch chan int) {
	for i := range 10 {
		ch <- i
	}
	close(ch) // Fecha o canal para que a função consume saiba que não haverá mais mensagens, evita deadlock
}

func consume(ch chan int) {
	for x := range ch {
		println("Received:", x)
	}
}
