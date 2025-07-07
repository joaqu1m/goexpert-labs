package main

func main() {

	channel := make(chan string)

	// thread 2
	go func() {
		channel <- "A" // isso faz a thread 2 esperar até que alguém consuma a mensagem
	}()

	// thread 1
	msg := <-channel // canal esvazia; caso ainda não tenha mensagem, a thread 1 espera até que alguém envie uma mensagem
	println("Received message:", msg)
}
