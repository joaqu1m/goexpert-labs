package main

import "time"

func task(name string) {
	for i := 0; i < 5; i++ {
		println(name, "is working on task", i)
		time.Sleep(time.Second)
	}
}

// main já é uma thread
// o garbage collector do Go também é uma green thread
func main() {

	// não utiliza concorrência
	// task("A")
	// task("B")

	// utiliza concorrência
	go task("A")
	go task("B")
	go func() {
		for i := range 5 {
			println("Anonymous task is working on task", i)
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(6 * time.Second)
}
