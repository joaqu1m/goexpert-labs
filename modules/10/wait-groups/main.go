package main

import (
	"sync"
	"time"
)

func task(name string, waitGroup *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		println(name, "is working on task", i)
		time.Sleep(time.Second)
	}
	waitGroup.Done() // sinaliza que a goroutine terminou
}

func main() {

	waitGroup := sync.WaitGroup{}
	// waitGroup com 3 créditos
	waitGroup.Add(3) // 1 crédito para cada goroutine

	go task("A", &waitGroup)
	go task("B", &waitGroup)
	go func() {
		for i := range 5 {
			println("Anonymous task is working on task", i)
			time.Sleep(time.Second)
		}
		waitGroup.Done() // sinaliza que a goroutine terminou
	}()

	waitGroup.Wait() // espera todas as goroutines terminarem
	println("All tasks completed")
}
