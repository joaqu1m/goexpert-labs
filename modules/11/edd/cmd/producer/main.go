package main

import (
	"github.com/joaqu1m/goexpert-labs/modules/11/edd/pkg"
)

func main() {
	ch, err := pkg.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	pkg.Publish(ch, []byte("Hello, RabbitMQ!"))
}
