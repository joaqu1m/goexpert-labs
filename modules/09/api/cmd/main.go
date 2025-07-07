package main

import (
	"fmt"

	"github.com/joaqu1m/goexpert-labs/modules/09/api/configs"
)

func main() {
	err := configs.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	fmt.Println(configs.Cfg.DBDriver)
}
