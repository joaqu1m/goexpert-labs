package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

type DBConfig struct {
	Host     string `json:"DB_HOST"`
	Port     string `json:"DB_PORT"`
	User     string `json:"DB_USER"`
	Password string `json:"DB_PASSWORD"`
	Name     string `json:"DB_NAME"`
}

var dbConfig DBConfig

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	UnmarshalConfig("env.json")

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					UnmarshalConfig("env.json")

					fmt.Println("Configuration reloaded from env.json")
					fmt.Println("DB Config:", dbConfig)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				panic(err)
			}
		}
	}()
	err = watcher.Add("env.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("Watching for changes in env.json...")
	<-done
	fmt.Println("Exiting...")
}

func UnmarshalConfig(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &dbConfig)
	if err != nil {
		panic(err)
	}
}
